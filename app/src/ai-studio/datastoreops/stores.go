package dmstores

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	datasvc "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices"
	models "github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
	utils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

// BeDbStore is the dbstores for Data Mesh data, containing the Redis client currently
type BeDataStore struct {
	Db                  *datasvc.DataStoreClient
	PubSub              *datasvc.DataStoreClient
	Cacher              *datasvc.DataStoreClient
	IsRedisStackEnabled bool
	Error               error
}

// func initDbStores(ctx context.Context, secName string, secretClient *BeSecretStore) (error) {
// }

// NewVapusBEDataStore creates a new BeDbStore object with driver client of different db backend store
func NewVapusBEDataStore(conf *utils.VapusAIStudioConfig, secretClient *BeSecretStore) *BeDataStore {
	ctx := context.Background()
	bds := &BeDataStore{}
	dbClient, err := initDbStores(ctx, conf.GetDBStoragePath(), secretClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing db data store")
	}
	bds.Db = dbClient
	cacheClient, err := initDbStores(ctx, conf.GetCachStoragePath(), secretClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing cache data store")
	}
	bds.Cacher = cacheClient

	return bds
}

func initDbStores(ctx context.Context, secName string, secretClient *BeSecretStore) (*datasvc.DataStoreClient, error) {
	bytes, err := secretClient.ReadSecret(ctx, secName)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while reading secret data for data store")
		return nil, err
	}
	creds := &models.ConnectionOpts{}
	log.Println("Secret data: ", string(bytes.([]byte)))
	err = json.Unmarshal(bytes.([]byte), creds)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while unmarshalling secret data")
		return nil, err
	}
	return datasvc.NewDataStoreClient(ctx, creds, logger)
}

func (ds *DMStore) GetDbStoreParams() *models.ConnectionOpts {
	return ds.BeDataStore.Db.ConnectionOpts
}

func (ds *DMStore) CreateIndex(ctx context.Context, index string) error {
	if exists, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Exists(index).Do(ctx); exists || err != nil {
		logger.Debug().Msgf("Index %v already exists.", index)
		return err
	}
	_, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Create(index).Do(ctx)
	if err != nil {
		logger.Fatal().Err(err).Ctx(ctx).Msg("error while creating index in elastic search")
		return err
	}
	logger.Debug().Msgf("Index %v created successfully.", index)
	return nil
}

func (ds *DMStore) CreateIndexWithMapping(ctx context.Context, index string, mapping *types.TypeMapping) error {
	if exists, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Exists(index).Do(ctx); exists || err != nil {
		logger.Debug().Msgf("Index %v already exists.", index)
		return err
	}
	_, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Create(index).Mappings(mapping).Do(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while updating index with mapping in elastic search, call error")
		return err
	}
	logger.Debug().Msgf("Index %v created successfully.", index)
	return nil
}
