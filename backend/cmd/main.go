package main

import (
	networkDomain "example.com/m/internal/domain/network"
	"example.com/m/internal/infrastructure/auth/oidc"
	"example.com/m/internal/infrastructure/persistence/sqlite"
	computeRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/compute"
	gatewayRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/gateway"
	imageRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/image"
	networkRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/network"
	storageRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/storage"
	userRepoInfra "example.com/m/internal/infrastructure/persistence/sqlite/user"
	"example.com/m/internal/interface/http"
	computeH "example.com/m/internal/interface/http/compute"
	userH "example.com/m/internal/interface/http/user"
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
	// uow := sqlite.NewUnitOfWork(db)
	userRepo := userRepoInfra.NewRepository(db)
	instRepo := computeRepoInfra.NewRepository(db)
	netRepo := networkRepoInfra.NewRepository(db)
	volRepo := storageRepoInfra.NewRepository(db)
	imgRepo := imageRepoInfra.NewRepository(db)
	ingressRepo := gatewayRepoInfra.NewRepository(db)

	// service
	netCalcuSvc := networkDomain.NewNetworkService()
	// identitySvc := userSvcInfra.NewService()

	// usecase
	// ensureUserUC := userUC.NewEnsureUserInteractor(userRepo, netRepo, /* identitySvc */ netCalcuSvc, /* publisher */, /* uow */)
	listMyInstanceUC := userUC.NewListMyInstanceInteractor(instRepo)
	// provisioningNetUC := networkUC.NewProvisioningNetworkInteractor(userRepo, netRepo, netCalcuSvc, /* identitySvc */, /* dvier */, /* uow */)
	// reqCreateUC := computeUC.NewRequestCreateInstanceInteractor(userRepo, instRepo, netRepo, ingressRepo, volRepo, netCalcuSvc, /* publisher */, /* uow */)
	// execCreateUC := computeUC.NewExecuteCreateInstanceInteractor(instRepo, volRepo, ingressRepo, /* instDriver */, /* volDriver */, /*ingressDriver */, /* uow */)
	
	// handler
	// userHandler := userH.NewHandler(oidcCfg, *oidcVerifier, /* ensureUserUC */, listMyInstanceUC)
	// computeHandler := computeH.NewHandler(/* reqCreateUC */, /* ensureUserUC */)

	// router
	// e := http.InitRoutes(/* userHandler */, /* computeHandler */)
	// e.Start(":8080")
}
