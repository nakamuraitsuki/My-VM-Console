package compute

import (
	"context"
	"log"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
)

type AuthorizePubKeyInput struct {
	InstanceID compute.InstanceID
	PubKey     string
}

type AuthorizePubkeyUseCase interface {
	Execute(ctx context.Context, input AuthorizePubKeyInput) error
}

type authorizePubKeyInteractor struct {
	instanceRepo   compute.InstanceRepository
	instanceDriver compute.ComputeDriver
}

func NewAuthorizePubKeyInteractor(
	instanceRepo compute.InstanceRepository,
	instanceDriver compute.ComputeDriver,
) AuthorizePubkeyUseCase {
	return &authorizePubKeyInteractor{
		instanceRepo:   instanceRepo,
		instanceDriver: instanceDriver,
	}
}

func (i *authorizePubKeyInteractor) Execute(ctx context.Context, input AuthorizePubKeyInput) error {
	// インスタンスの存在確認
	inst, err := i.instanceRepo.FindByID(ctx, input.InstanceID)
	if err != nil {
		return err
	}
	if inst == nil {
		return compute.ErrInstanceNotFound
	}

	// 権限検証
	usr, ok := user.FromContext(ctx)
	if !ok {
		return user.ErrUserNotInContext
	}
	// 見る権限がないし、所有者でもない場合は権限外
	if !usr.HasPermission(user.PermissionInstanceRead) && inst.OwnerID() != usr.ID() {
		return user.ErrNoPermission
	}

	// instanceの状態確認
	if inst.Status() != compute.StatusRunning {
		return compute.ErrInvalidInstanceStatus
	}

	// 公開鍵の登録
	if err := i.instanceDriver.AuthorizePublicKey(ctx, inst, input.PubKey); err != nil {
		log.Printf("[AuthorizePubKey] failed: AuthorizePublicKey instanceID=%s err=%v", input.InstanceID, err)
		return err
	}
	return nil
}
