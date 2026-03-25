package compute

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/gateway"
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
)

type RequestCreateInstanceInput struct {
	Name     string
	ImageID  image.ImageID
	VPCID    *network.VPCID    // Optional に対応したい
	SubnetID *network.SubnetID // Optional に対応したい
	CPU      int
	Memory   int
}

type RequestCreateInstanceOutput struct {
	InstanceID compute.InstanceID
	Name       string
	Status     compute.InstanceStatus
	SubnetID   network.SubnetID // 指定がなくても、割り当てたSubnetは特定出来る
	PrivateIP  string
}

type RequestCreateInstanceUseCase interface {
	Execute(ctx context.Context, req RequestCreateInstanceInput) (*RequestCreateInstanceOutput, error)
}

type requestCreateInstanceInteractor struct {
	userRepo       user.UserRepository
	instanceRepo   compute.InstanceRepository
	networkRepo    network.Repository
	gatewayRepo    gateway.Repository
	storageRepo    storage.Repository
	networkService network.NetworkService
	publisher      usecase.JobPublisher
	uow            usecase.UnitOfWork
}

func NewRequestCreateInstanceInteractor(
	userRepo user.UserRepository,
	instanceRepo compute.InstanceRepository,
	networkRepo network.Repository,
	gatewayRepo gateway.Repository,
	storageRepo storage.Repository,
	networkService network.NetworkService,
	publisher usecase.JobPublisher,
	uow usecase.UnitOfWork,
) RequestCreateInstanceUseCase {
	return &requestCreateInstanceInteractor{
		userRepo:       userRepo,
		instanceRepo:   instanceRepo,
		networkRepo:    networkRepo,
		gatewayRepo:    gatewayRepo,
		storageRepo:    storageRepo,
		networkService: networkService,
		publisher:      publisher,
		uow:            uow,
	}
}

func (i *requestCreateInstanceInteractor) Execute(
	ctx context.Context,
	req RequestCreateInstanceInput,
) (*RequestCreateInstanceOutput, error) {
	log.Printf("[RequestCreateInstance] start name=%s vpc=%v subnet=%v cpu=%d memory=%d", req.Name, req.VPCID, req.SubnetID, req.CPU, req.Memory)

	// check user permissions
	usr, ok := user.FromContext(ctx)
	if !ok {
		log.Printf("[RequestCreateInstance] failed: user not found in context")
		return nil, user.ErrUserNotInContext
	}
	log.Printf("[RequestCreateInstance] user resolved userID=%s", usr.ID())

	if !usr.HasPermission(user.PermissionInstanceCreate) {
		log.Printf("[RequestCreateInstance] failed: no permission userID=%s", usr.ID())
		return nil, user.ErrNoPermission
	}
	log.Printf("[RequestCreateInstance] permission check passed userID=%s", usr.ID())

	ownedInstances, err := i.instanceRepo.FindByOwnerID(ctx, usr.ID())
	if err != nil {
		log.Printf("[RequestCreateInstance] failed: FindByOwnerID userID=%s err=%v", usr.ID(), err)
		return nil, err
	}
	log.Printf("[RequestCreateInstance] owner instances loaded count=%d", len(ownedInstances))
	if !usr.CanAllocateInstance(len(ownedInstances), req.CPU) {
		log.Printf("[RequestCreateInstance] failed: quota exceeded userID=%s current=%d requestedCPU=%d", usr.ID(), len(ownedInstances), req.CPU)
		return nil, user.ErrQuotaExceeded
	}
	log.Printf("[RequestCreateInstance] quota check passed")

	var instanceID compute.InstanceID
	log.Printf("[RequestCreateInstance] starting unit of work")
	uowErr := i.uow.Do(ctx, func(ctx context.Context) error {
		// インスタンスを置くVPCを取得する
		if req.VPCID == nil {
			log.Printf("[RequestCreateInstance] loading vpc vpcID=%v", req.VPCID)
			vpc, err := i.networkRepo.FindVPCByUserID(ctx, usr.ID())
			if err != nil {
				log.Printf("[RequestCreateInstance] failed: FindVPCByID vpcID=%v err=%v", req.VPCID, err)
				return err
			}
			if vpc == nil {
				log.Printf("[RequestCreateInstance] failed: vpc not found vpcID=%v", req.VPCID)
				return network.ErrVPCNotFound
			}
			log.Printf("[RequestCreateInstance] vpc loaded vpcID=%s name=%s", vpc.ID(), vpc.Name())
			vpcID := vpc.ID() // ポインタを渡すために一旦変数に入れる
			req.VPCID = &vpcID // 以降はIDだけあればいいので,	IDを上書きしておく
		}

		// ついでにIPアドレスも確保する
		var targetSubnet *network.Subnet
		var IPAddress string
		var err error
		if req.SubnetID != nil {
			log.Printf("[RequestCreateInstance] subnet selection: explicit subnet=%s", *req.SubnetID)
			targetSubnet, err = i.networkRepo.FindSubnetByID(ctx, *req.SubnetID)
			if err != nil {
				log.Printf("[RequestCreateInstance] failed: FindSubnetByID subnet=%s err=%v", *req.SubnetID, err)
				return err
			}
			leases, err := i.networkRepo.FindLeasesBySubnetID(ctx, targetSubnet.ID())
			if err != nil {
				log.Printf("[RequestCreateInstance] failed: FindLeasesBySubnetID subnet=%s err=%v", targetSubnet.ID(), err)
				return err
			}
			log.Printf("[RequestCreateInstance] leases loaded subnet=%s count=%d", targetSubnet.ID(), len(leases))
			usedIPs := make([]string, len(leases))
			for i, lease := range leases {
				usedIPs[i] = lease.IPAddress
			}
			IPAddress, err = i.networkService.CalculateNextAvailableIP(ctx, targetSubnet.CIDR(), usedIPs)
			if err != nil {
				log.Printf("[RequestCreateInstance] failed: CalculateNextAvailableIP subnet=%s err=%v", targetSubnet.ID(), err)
				return err
			}
			log.Printf("[RequestCreateInstance] ip selected subnet=%s ip=%s", targetSubnet.ID(), IPAddress)
		} else {
			log.Printf("[RequestCreateInstance] subnet selection: auto by vpc=%v", req.VPCID)
			subnets, err := i.networkRepo.FindSubnetsByVPCID(ctx, *req.VPCID)
			if err != nil || len(subnets) == 0 {
				log.Printf("[RequestCreateInstance] failed: FindSubnetsByVPCID vpc=%s count=%d err=%v", *req.VPCID, len(subnets), err)
				return errors.New("no subnets available in the specified VPC")
			}
			log.Printf("[RequestCreateInstance] subnets loaded vpc=%s count=%d", *req.VPCID, len(subnets))

			var availableIP string
			for _, subnet := range subnets {
				log.Printf("[RequestCreateInstance] try subnet=%s cidr=%s", subnet.ID(), subnet.CIDR())
				leases, err := i.networkRepo.FindLeasesBySubnetID(ctx, subnet.ID())
				if err != nil {
					log.Printf("[RequestCreateInstance] failed: FindLeasesBySubnetID subnet=%s err=%v", subnet.ID(), err)
					return err
				}
				log.Printf("[RequestCreateInstance] leases loaded subnet=%s count=%d", subnet.ID(), len(leases))
				usedIPs := make([]string, len(leases))
				for i, lease := range leases {
					usedIPs[i] = lease.IPAddress
				}
				availableIP, err = i.networkService.CalculateNextAvailableIP(ctx, subnet.CIDR(), usedIPs)
				if err == nil {
					targetSubnet = subnet
					IPAddress = availableIP
					log.Printf("[RequestCreateInstance] ip selected subnet=%s ip=%s", subnet.ID(), IPAddress)
					break
				}
				// 「空きがない」以外のエラーは即座に返す
				if !errors.Is(err, network.ErrNoAvailableIPs) {
					log.Printf("[RequestCreateInstance] failed: CalculateNextAvailableIP subnet=%s err=%v", subnet.ID(), err)
					return err
				}
				log.Printf("[RequestCreateInstance] subnet exhausted subnet=%s", subnet.ID())
				// サブネットに利用可能なIPがない場合は次のサブネットを試す
			}
		}

		// サブネットが見つからない場合はエラー
		if targetSubnet == nil || IPAddress == "" {
			log.Printf("[RequestCreateInstance] failed: no subnet/ip available")
			return errors.New("no available subnets with free IP addresses")
		}

		instanceID = compute.NewID()
		log.Printf("[RequestCreateInstance] generated instance id=%s", instanceID)
		lease := network.NewLease(
			targetSubnet.ID(),
			string(instanceID),
			IPAddress,
		)
		log.Printf("[Request CreateInstance] Reserving IP %s for instance %s in subnet %s", IPAddress, instanceID, targetSubnet.ID())
		if err := i.networkRepo.CreateLease(ctx, lease); err != nil {
			log.Printf("[RequestCreateInstance] failed: CreateLease subnet=%s ip=%s err=%v", targetSubnet.ID(), IPAddress, err)
			return err
		}
		log.Printf("[RequestCreateInstance] lease created subnet=%s ip=%s", targetSubnet.ID(), IPAddress)

		// volume 予約
		volumeID := storage.NewID()
		defaultSize := 20 // GB
		volume := storage.NewVolume(
			volumeID,
			req.Name+"-root",
			defaultSize,
			"zfs",            // 仮のプール名
			string(usr.ID()), // とりあえずユーザーの持ち物にしておく
		)
		if err := i.storageRepo.Save(ctx, volume); err != nil {
			log.Printf("[RequestCreateInstance] failed: Save volume volumeID=%s err=%v", volumeID, err)
			return err
		}
		log.Printf("[RequestCreateInstance] volume saved volumeID=%s", volumeID)

		// Save Entity
		inst := compute.NewInstance(
			instanceID,
			req.Name,
			usr.ID(),
			compute.StatusPending, // 初期状態はPending
			nil,
			req.CPU,
			req.Memory,
			req.ImageID,
			*req.VPCID,
			targetSubnet.ID(),
			IPAddress,
			volume.ID(),
		)

		if err := i.instanceRepo.Save(ctx, inst); err != nil {
			log.Printf("[RequestCreateInstance] failed: Save instance instance=%v err=%v", inst, err)
			return err
		}
		log.Printf("[RequestCreateInstance] instance saved instanceID=%s", instanceID)

		// インスタンスを外部に公開する
		ingressID := gateway.NewID()
		subdomain := req.Name
		// 初期に与えるのはSSH用ルートのみ。必要になったらユーザー自身が追加する。
		ingressRoute := gateway.NewIngressRoute(
			ingressID,
			subdomain,
			"ssh",
			IPAddress,
			22, // 仮のポート番号
			string(usr.ID()),
			instanceID,
		)
		if err := i.gatewayRepo.Save(ctx, ingressRoute); err != nil {
			log.Printf("[RequestCreateInstance] failed: Save ingress route ingressID=%s instanceID=%s err=%v", ingressID, instanceID, err)
			return err
		}
		log.Printf("[RequestCreateInstance] ingress route saved ingressID=%s instanceID=%s", ingressID, instanceID)

		return nil
	})
	if uowErr != nil {
		log.Printf("[RequestCreateInstance] failed: uow err=%v", uowErr)
		return nil, uowErr
	}
	log.Printf("[RequestCreateInstance] uow committed instanceID=%s", instanceID)

	// ジョブの発行
	payload := CreateInstancePayload{
		InstanceID: instanceID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[RequestCreateInstance] failed: marshal payload instanceID=%s err=%v", instanceID, err)
		return nil, err
	}
	log.Printf("[RequestCreateInstance] job payload marshaled instanceID=%s", instanceID)
	if err := i.publisher.Publish(ctx, usecase.JobTypeCreateInstance, payloadBytes); err != nil {
		log.Printf("[RequestCreateInstance] failed: publish job instanceID=%s err=%v", instanceID, err)
		inst, err := i.instanceRepo.FindByID(ctx, instanceID)
		if err != nil {
			log.Printf("[RequestCreateInstance] failed: FindByID after publish failure instanceID=%s err=%v", instanceID, err)
			return nil, err
		}
		inst.MarkAsError(compute.ErrInPending)
		if saveErr := i.instanceRepo.Save(ctx, inst); saveErr != nil {
			log.Printf("[RequestCreateInstance] failed: Save error state instanceID=%s err=%v", instanceID, saveErr)
			return nil, saveErr
		}
		log.Printf("[RequestCreateInstance] marked instance error state instanceID=%s", instanceID)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return nil, err
	}
	log.Printf("[RequestCreateInstance] job published instanceID=%s", instanceID)

	createdInstance, err := i.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		log.Printf("[RequestCreateInstance] failed: FindByID final instanceID=%s err=%v", instanceID, err)
		return nil, err
	}
	log.Printf("[RequestCreateInstance] completed instanceID=%s status=%s subnet=%s ip=%s", createdInstance.ID(), createdInstance.Status(), createdInstance.SubnetID(), createdInstance.PrivateIP())
	return &RequestCreateInstanceOutput{
		InstanceID: createdInstance.ID(),
		Name:       createdInstance.Name(),
		Status:     createdInstance.Status(),
		SubnetID:   createdInstance.SubnetID(),
		PrivateIP:  createdInstance.PrivateIP(),
	}, nil
}
