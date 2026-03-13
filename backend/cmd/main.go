package main

import (
	"context"
	"net/http"
	"time"

	networkDomain "example.com/m/internal/domain/network"
	authHttp "example.com/m/internal/infrastructure/auth/http"
	"example.com/m/internal/infrastructure/auth/mapper"
	"example.com/m/internal/infrastructure/auth/oidc"
	"example.com/m/internal/infrastructure/driver/etcd"
	gatewayDriverInfra "example.com/m/internal/infrastructure/driver/etcd/gateway"
	"example.com/m/internal/infrastructure/driver/incus"
	computeDriverInfra "example.com/m/internal/infrastructure/driver/incus/compute"
	networkDriverInfra "example.com/m/internal/infrastructure/driver/incus/network"
	storageDriverInfra "example.com/m/internal/infrastructure/driver/incus/storage"
	"example.com/m/internal/infrastructure/env"
	"example.com/m/internal/infrastructure/messaging/channel"
	"example.com/m/internal/infrastructure/persistence/sqlite"
	computeRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/compute"
	gatewayRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/gateway"
	imageRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/image"
	networkRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/network"
	storageRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/storage"
	userRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/user"
	httpAdapter "example.com/m/internal/interface/http"
	computeH "example.com/m/internal/interface/http/compute"
	userH "example.com/m/internal/interface/http/user"
	"example.com/m/internal/usecase"
	computeUC "example.com/m/internal/usecase/compute"
	networkUC "example.com/m/internal/usecase/network"
	userUC "example.com/m/internal/usecase/user"
)

func main() {

	// DB client
	dbCfg := sqlite.NewConfig()
	db := sqlite.NewSqlxClient(dbCfg)

	// OIDC
	oidcCfg := oidc.NewOIDCConfig()
	oidcVerifier := oidc.NewIDTokenVerifier(oidcCfg)

	// repository
	uow := sqlite.NewTransactor(db)
	userRepo := userRepoInfra.NewRepository(db)
	instRepo := computeRepoInfra.NewRepository(db)
	netRepo := networkRepoInfra.NewRepository(db)
	volRepo := storageRepoInfra.NewRepository(db)
	imgRepo := imageRepoInfra.NewRepository(db)
	ingressRepo := gatewayRepoInfra.NewRepository(db)

	// service
	// TODO: 権限のマッピングのためのJSONを用意すること
	filePath := env.GetString("EXTERNAL_AUTH_PERMISSION_MAPPING_FILE", "config/permission_mapping.yaml")
	apiEndpoint := env.GetString("OAUTH_API_ENDPOINT", "")
	permissionMapper := mapper.NewPermissionMapper(filePath, 3*time.Hour)
	netCalcuSvc := networkDomain.NewNetworkService()
	identitySvc := authHttp.NewExternalIdentityService(http.DefaultClient, apiEndpoint, permissionMapper)

	// driver
	incusConfig := incus.NewConfig()
	incusClient := incus.NewClient(incusConfig)
	instDriver := computeDriverInfra.NewDriver(incusClient)
	netDriver := networkDriverInfra.NewDriver(incusClient)
	volDriver := storageDriverInfra.NewDriver(incusClient)

	etcdConfig := etcd.NewConfig()
	etcdClient := etcd.NewClient(etcdConfig)
	ingressDriver := gatewayDriverInfra.NewDriver(etcdClient)

	// worker
	hub := channel.NewHub()
	publisher := channel.NewPublisher(hub)
	subscriber := channel.NewSubscriber(hub)

	// usecase
	ensureUserUC := userUC.NewEnsureUserInteractor(userRepo, netRepo, identitySvc, netCalcuSvc, publisher, uow)
	listMyInstanceUC := userUC.NewListMyInstanceInteractor(instRepo)
	provisioningNetUC := networkUC.NewProvisioningNetworkInteractor(userRepo, netRepo, netCalcuSvc, identitySvc, netDriver, uow)
	reqCreateUC := computeUC.NewRequestCreateInstanceInteractor(userRepo, instRepo, netRepo, ingressRepo, volRepo, netCalcuSvc, publisher, uow)
	execCreateUC := computeUC.NewExecuteCreateInstanceInteractor(instRepo ,netRepo, volRepo, ingressRepo, imgRepo, instDriver, volDriver, ingressDriver, uow)

	// handler
	userHandler := userH.NewHandler(oidcCfg, *oidcVerifier, ensureUserUC, listMyInstanceUC)
	computeHandler := computeH.NewHandler(reqCreateUC, ensureUserUC)

	// bind job handlers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	usecase.Bind(ctx, subscriber, usecase.JobTypeCreateVPCAndDefaultSubnet, provisioningNetUC.Execute)
	usecase.Bind(ctx, subscriber, usecase.JobTypeCreateInstance, execCreateUC.Execute)
	// router
	e := httpAdapter.InitRoutes(userHandler, computeHandler)
	e.Start(":8080")
}
