package dmstores

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	datasvc "github.com/vapusdata-ecosystem/vapusai-studio/internals/dataservices"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
	"gopkg.in/yaml.v3"
)

type BeSecretStore struct {
	SecretClient BeSecretStoreOps
	Error        error
}

type BeSecretStoreOps interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
}

func LoadSecretStoreParams(filePath string, log zerolog.Logger) (*models.ConnectionOpts, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Msgf("Error reading file: %v", err)
		return nil, err
	}
	conf := &models.ConnectionOpts{}
	log.Info().Msgf("File Content after reading : %s", string(bytes))
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		log.Fatal().Msgf("Error unmarshalling file: %v", err)
		return nil, err
	}
	err = validator.New().Struct(conf)

	if err != nil {
		return nil, err
	}
	return conf, nil
}

// NewVapusBESecretStore creates a new BeSecretStore object with driver client of different secret backend store based on the configuration provided
func NewVapusBESecretStore(path string) (*BeSecretStore, error) {
	logger.Debug().Msgf("Creating secret store client with path: %s", path)
	ctx := context.Background()
	connOpts, err := LoadSecretStoreParams(path, logger)
	if err != nil {
		logger.Info().Msgf("Error while loading secret store params: %v", err)
		return nil, err
	}
	client, err := datasvc.NewSecretStoreClient(ctx, connOpts)
	if err != nil {
		return nil, err
	}
	return &BeSecretStore{
		SecretClient: client,
	}, nil
}

// WriteSecret writes the secret to the secret store
func (be *BeSecretStore) WriteSecret(ctx context.Context, name string, secrdData any) error {
	logger.Info().Msgf("Writing secret %v to secret store", name)
	switch secrdData.(type) { //nolint:all
	case map[string]interface{}:
		return be.SecretClient.WriteSecret(ctx, secrdData.(map[string]interface{}), name) //nolint:all
	case string:
		return be.SecretClient.WriteSecret(ctx, map[string]interface{}{"value": secrdData.(string)}, name) //nolint:all
	case []string:
		return be.SecretClient.WriteSecret(ctx, map[string]interface{}{"value": secrdData.(string)}, name)
	case []interface{}:
		return be.SecretClient.WriteSecret(ctx, map[string]interface{}{"value": secrdData.(string)}, name)
	default:
		// handle other data types like slice
		// return be.Client.WriteSecret(ctx, map[string]interface{}{"value": string(secrdData.(interface{}))}, name)
		return be.SecretClient.WriteSecret(ctx, map[string]interface{}{"value": secrdData.(string)}, name)
	}

}

// ReadSecret reads the secret from the secret store
func (be *BeSecretStore) ReadSecret(ctx context.Context, secretId string) (any, error) {
	return be.SecretClient.ReadSecret(ctx, secretId)
}

// DeleteSecret deletes the secret from the secret store
func (be *BeSecretStore) DeleteSecret(ctx context.Context, secretId string) error {
	return be.SecretClient.DeleteSecret(ctx, secretId)
}
