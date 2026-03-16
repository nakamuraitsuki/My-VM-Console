package network

import (
	"context"
	"log"

	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
)

type CreateVPCAndDefaultSubnetPayload struct {
	VPCID    network.VPCID    `json:"vpc_id"`
	SubnetID network.SubnetID `json:"subnet_id"`
	Identity struct {
		DisplayName     string   `json:"display_name"`
		ProfileImageURL string   `json:"profile_image_url"`
		Permissions     []string `json:"permissions"`
	} `json:"identity"`
}

type ProvisioningNetworkUseCase interface {
	Execute(ctx context.Context, payload CreateVPCAndDefaultSubnetPayload) error
}

type provisioningNetworkInteractor struct {
	userRepo    user.UserRepository
	networkRepo network.Repository
	networkSvc  network.NetworkService
	identitySvc user.IdentityService
	driver      network.NetworkDriver
	uow         usecase.UnitOfWork
}

func NewProvisioningNetworkInteractor(
	userRepo user.UserRepository,
	networkRepo network.Repository,
	networkSvc network.NetworkService,
	identitySvc user.IdentityService,
	driver network.NetworkDriver,
	uow usecase.UnitOfWork,
) ProvisioningNetworkUseCase {
	return &provisioningNetworkInteractor{
		userRepo:    userRepo,
		networkRepo: networkRepo,
		networkSvc:  networkSvc,
		identitySvc: identitySvc,
		driver:      driver,
		uow:         uow,
	}
}

// Jobが投げられた時点で、User作成とVPC、Subnetの予約は完了している前提
func (i *provisioningNetworkInteractor) Execute(ctx context.Context, payload CreateVPCAndDefaultSubnetPayload) error {
	log.Printf("[ProvisioningNetwork] start: vpc_id=%s subnet_id=%s", payload.VPCID, payload.SubnetID)

	log.Printf("[ProvisioningNetwork] step=find_vpc_by_id vpc_id=%s", payload.VPCID)
	vpc, err := i.networkRepo.FindVPCByID(ctx, payload.VPCID)
	if err != nil {
		log.Printf("[ProvisioningNetwork] step=find_vpc_by_id result=error vpc_id=%s err=%v", payload.VPCID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=find_vpc_by_id result=success vpc_id=%s", payload.VPCID)

	log.Printf("[ProvisioningNetwork] step=find_subnet_by_id subnet_id=%s", payload.SubnetID)
	subnet, err := i.networkRepo.FindSubnetByID(ctx, payload.SubnetID)
	if err != nil {
		log.Printf("[ProvisioningNetwork] step=find_subnet_by_id result=error subnet_id=%s err=%v", payload.SubnetID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=find_subnet_by_id result=success subnet_id=%s", payload.SubnetID)

	userID := user.UserID(vpc.OwnerID())
	log.Printf("[ProvisioningNetwork] step=find_user_by_id user_id=%s", userID)
	userData, err := i.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("[ProvisioningNetwork] step=find_user_by_id result=error user_id=%s err=%v", userID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=find_user_by_id result=success user_id=%s", userID)
	permissions := make([]user.Permission, 0, len(payload.Identity.Permissions))
	for _, permStr := range payload.Identity.Permissions {
		permissions = append(permissions, user.Permission(permStr))
	}
	usr := user.NewUser(
		userID,
		payload.Identity.DisplayName,
		payload.Identity.ProfileImageURL,
		permissions,
		userData.Quota,
		userData.Status,
		userData.ErrorPhase,
	)

	log.Printf("[ProvisioningNetwork] step=mark_user_initializing user_id=%s", userID)
	if err := usr.MarkAsInitializing(); err != nil {
		log.Printf("[ProvisioningNetwork] step=mark_user_initializing result=error user_id=%s err=%v", userID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=mark_user_initializing result=success user_id=%s", userID)

	log.Printf("[ProvisioningNetwork] step=save_user_initializing user_id=%s", userID)
	if err := i.userRepo.Save(ctx, usr); err != nil {
		log.Printf("[ProvisioningNetwork] step=save_user_initializing result=error user_id=%s err=%v", userID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=save_user_initializing result=success user_id=%s", userID)

	log.Printf("[ProvisioningNetwork] step=create_vpc_driver vpc_id=%s", vpc.ID())
	if err := i.driver.CreateVPC(ctx, vpc); err != nil {
		log.Printf("[ProvisioningNetwork] step=create_vpc_driver result=error vpc_id=%s err=%v", vpc.ID(), err)
		log.Printf("[ProvisioningNetwork] step=mark_user_failed_on_vpc_error user_id=%s", userID)
		usr.MarkAsFailed(user.FailedInInitializing)
		if saveErr := i.userRepo.Save(ctx, usr); saveErr != nil {
			log.Printf("[ProvisioningNetwork] step=save_user_failed_on_vpc_error result=error user_id=%s err=%v", userID, saveErr)
		} else {
			log.Printf("[ProvisioningNetwork] step=save_user_failed_on_vpc_error result=success user_id=%s", userID)
		}
		return err
	}
	log.Printf("[ProvisioningNetwork] step=create_vpc_driver result=success vpc_id=%s", vpc.ID())

	log.Printf("[ProvisioningNetwork] step=create_subnet_driver subnet_id=%s vpc_id=%s", subnet.ID(), vpc.ID())
	if err := i.driver.CreateSubnet(ctx, vpc.ID(), subnet); err != nil {
		log.Printf("[ProvisioningNetwork] step=create_subnet_driver result=error subnet_id=%s vpc_id=%s err=%v", subnet.ID(), vpc.ID(), err)
		log.Printf("[ProvisioningNetwork] step=mark_user_failed_on_subnet_error user_id=%s", userID)
		usr.MarkAsFailed(user.FailedInInitializing)
		if saveErr := i.userRepo.Save(ctx, usr); saveErr != nil {
			log.Printf("[ProvisioningNetwork] step=save_user_failed_on_subnet_error result=error user_id=%s err=%v", userID, saveErr)
		} else {
			log.Printf("[ProvisioningNetwork] step=save_user_failed_on_subnet_error result=success user_id=%s", userID)
		}
		return err
	}
	log.Printf("[ProvisioningNetwork] step=create_subnet_driver result=success subnet_id=%s vpc_id=%s", subnet.ID(), vpc.ID())

	log.Printf("[ProvisioningNetwork] step=mark_user_active user_id=%s", userID)
	if err := usr.MarkAsActive(); err != nil {
		log.Printf("[ProvisioningNetwork] step=mark_user_active result=error user_id=%s err=%v", userID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=mark_user_active result=success user_id=%s", userID)

	log.Printf("[ProvisioningNetwork] step=save_user_active user_id=%s", userID)
	if err := i.userRepo.Save(ctx, usr); err != nil {
		log.Printf("[ProvisioningNetwork] step=save_user_active result=error user_id=%s err=%v", userID, err)
		return err
	}
	log.Printf("[ProvisioningNetwork] step=save_user_active result=success user_id=%s", userID)
	log.Printf("[ProvisioningNetwork] done: vpc_id=%s subnet_id=%s user_id=%s", payload.VPCID, payload.SubnetID, userID)
	return nil
}
