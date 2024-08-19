package server

import (
	"context"
	"path/filepath"
	"sync"

	config "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/config"
	dmstores "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/services"
	datasvc "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices"
	es "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices/elasticsearch"
	utils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

func packagesInit() {
	//Initialize the logger
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "vapusPlatform server init")

	logger.Info().Msg("Loading service config for VapusData server")
	// Load the service configuration, secrets inton the memory of the service. These information will be used by the service to connect to the database, vault etc connections
	config.InitServiceConfig(flagconfPath, filepath.Join(flagconfPath, configName))

	logger.Info().Msg("Service config loaded successfully")

	ctx := context.Background()

	bootStores(ctx, config.ServiceConfigManager)
	logger.Info().Msg("Service data stores loaded successfully")

	logger.Info().Msg("Service store dependencies loaded successfully")

	logger.Info().Msg("Service config loaded successfully")
	// Initialize the jwt authn validator
	logger.Info().Msgf("Loading JWT authn with secret path: %s", config.ServiceConfigManager.GetJwtAuthSecretPath())
	// Initialize the authn manager
	logger.Info().Msgf("Loading Authn manager with secret path: %v", config.AuthnParams)
	logger.Info().Msgf("Loading VapusAuth with secret path - %v", config.JwtAuthnParams)
	pkgs.InitAuthnManager(config.AuthnParams)
	// Initialize the NewVapusAuth
	pkgs.InitAuthzManager(config.JwtAuthnParams, pkgs.DmValidator)

	defer ctx.Done()
}

func bootStores(ctx context.Context, conf *utils.VapusAIStudioConfig) {
	//Boot the stores
	logger.Info().Msg("Booting the data stores")
	dmstores.InitDMStore(conf)
	if dmstores.DMStoreManager.Error != nil {
		logger.Fatal().Err(dmstores.DMStoreManager.Error).Msg("error while initializing data stores.")
	}
	services.InitDMServices(dmstores.DMStoreManager)
	bootEsIndex(ctx, dmstores.DMStoreManager.Db.ElasticSearchClient)
}

func bootEsIndex(ctx context.Context, cl *es.ElasticSearchStore) {
	if cl == nil {
		logger.Fatal().Msg("error while booting ES indexes")
	}
	errChan := make(chan error, len(datasvc.INDEX_LIST))
	var wg sync.WaitGroup
	for _, index := range datasvc.INDEX_LIST {
		wg.Add(1)
		go func(index string, wg *sync.WaitGroup) {
			defer wg.Done()
			val, ok := datasvc.INDEX_MAPPING[index]
			if ok {
				err := dmstores.DMStoreManager.CreateIndexWithMapping(ctx, index, val)
				if err != nil {
					errChan <- err
				}
			} else {
				err := dmstores.DMStoreManager.CreateIndex(ctx, index)
				if err != nil {
					errChan <- err
				}
			}
		}(index, &wg)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			logger.Err(err).Msg("error while booting ES index ")
		}
	}
}
