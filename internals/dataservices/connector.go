package datasvc

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	velasticsearch "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices/elasticsearch"
	ves "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices/elasticsearch"
	vredis "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices/redis"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"

	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
	dmerrors "github.com/vapusdata-ecosystem/vapusai-studio/internals/errors"
	logger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
	tpaws "github.com/vapusdata-ecosystem/vapusai-studio/internals/thirdparty/aws"
	tpazure "github.com/vapusdata-ecosystem/vapusai-studio/internals/thirdparty/azure"
	tpgcp "github.com/vapusdata-ecosystem/vapusai-studio/internals/thirdparty/gcp"
	tphcvault "github.com/vapusdata-ecosystem/vapusai-studio/internals/thirdparty/hcvault"
)

type DataStoreClient struct {
	RedisClient         *vredis.RedisStore
	ElasticSearchClient *velasticsearch.ElasticSearchStore
	ConnectionOpts      *models.ConnectionOpts
}

type SecretStoreOps interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
}

type SecretStoreClient struct {
	SecretStoreOps
}

var (
	log = logger.CoreLogger
)

func (d *DataStoreClient) Close() {
	if d.RedisClient != nil {
		d.RedisClient.Close()
	}
}

func (d *SecretStoreClient) Close() {

}

// NewDataConnClient creates a new DataConnClient
func NewDataStoreClient(ctx context.Context, opts *models.ConnectionOpts, log zerolog.Logger) (*DataStoreClient, error) {
	log.Debug().Msgf("Creating new data source connection for %s", opts.DataSourceEngine)
	switch opts.DataSourceEngine {
	case mpb.StorageEngine_REDIS.String():
		client, err := connectRedis(ctx, opts)
		if err != nil {
			log.Err(err).Msg("Error connecting to redis")
			return nil, err
		}
		return &DataStoreClient{RedisClient: client, ConnectionOpts: opts}, nil
	case mpb.StorageEngine_ELASTICSEARCH.String():
		client, err := connectElaticSearch(ctx, opts)
		if err != nil {
			log.Err(err).Msg("Error connecting to elasticsearch")
			return nil, err
		}
		return &DataStoreClient{ElasticSearchClient: client, ConnectionOpts: opts}, nil
	default:
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
}

func connectRedis(ctx context.Context, opts *models.ConnectionOpts) (*vredis.RedisStore, error) {
	return vredis.NewRedisStore(ctx, &vredis.Redis{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Password: opts.Credentials.Password,
		Username: opts.Credentials.Username,
	}, log)
}

func connectElaticSearch(ctx context.Context, opts *models.ConnectionOpts) (*ves.ElasticSearchStore, error) {
	return ves.NewElasticSearchStore(&ves.ElasticSearch{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Username: opts.Credentials.Username,
		Password: opts.Credentials.Password,
		ApiKey:   opts.Credentials.ApiToken,
	}, log)
}

func NewSecretStoreClient(ctx context.Context, opts *models.ConnectionOpts) (*SecretStoreClient, error) {
	log.Debug().Msg("Creating new secret store client")

	switch opts.DataSourceEngine {
	case mpb.StorageEngine_HASHICORPVAULT.String():
		var se string
		if val, ok := opts.Params[SECRETENGINE]; ok {
			se = val.(string)
		}
		client, err := tphcvault.NewHcVaultManager(ctx, &tphcvault.Vault{
			URL: opts.URL,
			// AuthAppRole:     conf.HashicorpVault.AppRoleAuthnEnabled,
			// ApproleRoleID:   conf.HashicorpVault.AppRoleID,
			// ApproleSecretID: conf.HashicorpVault.AppRoleSecret,
			Token:        opts.Credentials.GetApiToken(),
			SecretEngine: se,
		})
		if err != nil {
			log.Err(err).Msg("Error creating vault client")
			return &SecretStoreClient{}, err
		}
		return &SecretStoreClient{SecretStoreOps: client}, nil
	case mpb.StorageEngine_AWS_VAULT.String():
		client, err := tpaws.NewAwsSmClient(ctx, &tpaws.AWSConfig{
			Region:          opts.Credentials.AwsCreds.Region,
			AccessKeyId:     opts.Credentials.AwsCreds.AccessKeyId,
			SecretAccessKey: opts.Credentials.AwsCreds.SecretAccessKey,
		})
		if err != nil {
			log.Err(err).Msg("Error creating aws secret manager client")
			return &SecretStoreClient{}, err
		}
		return &SecretStoreClient{SecretStoreOps: client}, nil
	case mpb.StorageEngine_GCP_VAULT.String():
		decodeData, err := base64.StdEncoding.DecodeString(opts.Credentials.GcpCreds.ServiceAccountKey)
		if err != nil {
			log.Err(err).Msg("Error decoding gcp service account key")
			return &SecretStoreClient{}, err
		}
		client, err := tpgcp.NewGcpSMStore(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         opts.Credentials.GcpCreds.ProjectId,
			Region:            opts.Credentials.GcpCreds.Region,
		})
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return &SecretStoreClient{}, err
		}
		return &SecretStoreClient{SecretStoreOps: client}, nil
	case mpb.StorageEngine_AZURE_VAULT.String():
		client, err := tpazure.NewAzureKeyVault(ctx, &tpazure.AzureConfig{
			TenantID:     opts.Credentials.AzureCreds.TenantId,
			ClientID:     opts.Credentials.AzureCreds.ClientId,
			ClientSecret: opts.Credentials.AzureCreds.ClientSecret,
		}, opts.URL)
		if err != nil {
			log.Err(err).Msg("Error creating azure key vault client")
			return &SecretStoreClient{}, err
		}
		return &SecretStoreClient{SecretStoreOps: client}, nil
	default:
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
}
