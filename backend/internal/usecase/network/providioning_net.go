package network

import (
	"context"

	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
)

type CreateVPCAndDefaultSubnetPayload struct {
	VPCID    network.VPCID    `json:"vpc_id"`
	SubnetID network.SubnetID `json:"subnet_id"`
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
	vpc, err := i.networkRepo.FindVPCByID(ctx, payload.VPCID)
	if err != nil {
		return err
	}
	subnet, err := i.networkRepo.FindSubnetByID(ctx, payload.SubnetID)
	if err != nil {
		return err
	}

	userID := user.UserID(vpc.OwnerID())
	userData, err := i.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	identity, err := i.identitySvc.GetIdentity(ctx, string(userID))
	if err != nil {
		return err
	}

	usr := user.NewUser(
		userID,
		identity.DisplayName,
		identity.Permissions,
		userData.Quota,
		userData.Status,
		userData.ErrorPhase,
	)

	if err := usr.MarkAsInitializing(); err != nil {
		return err
	}
	if err := i.userRepo.Save(ctx, usr); err != nil {
		return err
	}

	if err := i.driver.CreateVPC(ctx, vpc); err != nil {
		usr.MarkAsFailed(user.FailedInInitializing)
		_ = i.userRepo.Save(ctx, usr) // エラー状態を保存
		return err
	}

	if err := i.driver.CreateSubnet(ctx, vpc.ID(), subnet); err != nil {
		usr.MarkAsFailed(user.FailedInInitializing)
		_ = i.userRepo.Save(ctx, usr) // エラー状態を保存
		return err
	}

	if err := usr.MarkAsActive(); err != nil {
		return err
	}
	if err := i.userRepo.Save(ctx, usr); err != nil {
		return err
	}
	return nil
}
