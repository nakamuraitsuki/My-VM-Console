package compute

import (
	"context"
	"errors"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
	"github.com/google/uuid"
)

type RequestCreateInstanceInput struct {
	Name     string
	OwnerID  user.UserID
	ImageID  image.ImageID
	VPCID    network.VPCID
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
	instanceRepo   compute.InstanceRepository
	networkRepo    network.Repository
	storageRepo    storage.Repository
	networkService network.NetworkService
	publisher      usecase.JobPublisher
	uow            usecase.UnitOfWork
}

func NewRequestCreateInstanceInteractor(
	instanceRepo compute.InstanceRepository,
	networkRepo network.Repository,
	storageRepo storage.Repository,
	networkService network.NetworkService,
	publisher usecase.JobPublisher,
	uow usecase.UnitOfWork,
) RequestCreateInstanceUseCase {
	return &requestCreateInstanceInteractor{
		instanceRepo:   instanceRepo,
		networkRepo:    networkRepo,
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
	var instanceID compute.InstanceID
	uowErr := i.uow.Do(ctx, func(ctx context.Context) error {
		// インスタンスを置くサブネットを取得する
		// ついでにIPアドレスも確保する
		var targetSubnet *network.Subnet
		var IPAddress string
		var err error
		if req.SubnetID != nil {
			targetSubnet, err = i.networkRepo.FindSubnetByID(ctx, *req.SubnetID)
			if err != nil {
				return err
			}
			leases, err := i.networkRepo.FindLeasesBySubnetID(ctx, targetSubnet.ID())
			if err != nil {
				return err
			}
			usedIPs := make([]string, len(leases))
			for i, lease := range leases {
				usedIPs[i] = lease.IPAddress
			}
			IPAddress, err = i.networkService.CalculateNextAvailableIP(ctx, targetSubnet.CIDR(), usedIPs)
		} else {
			subnets, err := i.networkRepo.FindSubnetsByVPCID(ctx, req.VPCID)
			if err != nil || len(subnets) == 0 {
				return errors.New("no subnets available in the specified VPC")
			}

			var availableIP string
			for _, subnet := range subnets {
				leases, err := i.networkRepo.FindLeasesBySubnetID(ctx, subnet.ID())
				if err != nil {
					return err
				}
				usedIPs := make([]string, len(leases))
				for i, lease := range leases {
					usedIPs[i] = lease.IPAddress
				}
				availableIP, err = i.networkService.CalculateNextAvailableIP(ctx, subnet.CIDR(), usedIPs)
				if err == nil {
					targetSubnet = subnet
					IPAddress = availableIP
					break
				}
				// 「空きがない」以外のエラーは即座に返す
				if !errors.Is(err, network.ErrNoAvailableIPs) {
					return err
				}
				// サブネットに利用可能なIPがない場合は次のサブネットを試す
			}
		}

		// サブネットが見つからない場合はエラー
		if targetSubnet == nil || IPAddress == "" {
			return errors.New("no available subnets with free IP addresses")
		}

		instanceID = compute.InstanceID("instance-" + uuid.New().String())
		lease := network.NewLease(
			targetSubnet.ID(),
			string(instanceID),
			IPAddress,
		)
		if err := i.networkRepo.CreateLease(ctx, lease); err != nil {
			return err
		}

		// volume 予約
		volumeID := storage.VolumeID("volume-" + uuid.New().String())
		defaultSize := 20 // GB
		volume := storage.NewVolume(
			volumeID,
			req.Name+"-root",
			defaultSize,
			"zfs", // 仮のプール名
			string(req.OwnerID),
		)
		if err := i.storageRepo.Save(ctx, volume); err != nil {
			return err
		}

		// Save Entity
		inst := compute.NewInstance(
			instanceID,
			req.Name,
			req.OwnerID,
			compute.StatusPending, // 初期状態はPending
			req.CPU,
			req.Memory,
			req.ImageID,
			targetSubnet.ID(),
			IPAddress,
			volume.ID(),
		)

		if err := i.instanceRepo.Save(ctx, inst); err != nil {
			return err
		}

		// ジョブの発行
		payload := CreateInstancePayload{
			InstanceID: instanceID,
		}
		if err := i.publisher.Publish(ctx, usecase.JobTypeCreateInstance, payload); err != nil {
			return err
		}
		return nil
	})
	if uowErr != nil {
		return nil, uowErr
	}

	createdInstance, err := i.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	return &RequestCreateInstanceOutput{
		InstanceID: createdInstance.ID(),
		Name:       createdInstance.Name(),
		Status:     createdInstance.Status(),
		SubnetID:   createdInstance.SubnetID(),
		PrivateIP:  createdInstance.PrivateIP(),
	}, nil
}
