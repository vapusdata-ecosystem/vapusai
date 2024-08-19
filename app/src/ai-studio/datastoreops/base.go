package dmstores

import (
	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

type DMStore struct {
	// Inheriting VapusBESecretStorage for secret store
	*BeSecretStore

	// Backend store for data mesh data
	*BeDataStore

	Error error
}

// GLobal var for DM store, it can accessed across the service
var (
	DMStoreManager *DMStore
	logger         zerolog.Logger
)

// Constructor to create new object for DMStore struct
func newDMStore(conf *utils.VapusAIStudioConfig) *DMStore {
	dmSec, err := NewVapusBESecretStore(conf.GetSecretStoragePath())
	if err != nil {
		return &DMStore{
			Error: err,
		}
	}
	dmDb := NewVapusBEDataStore(conf, dmSec)
	if dmDb.Error != nil {
		return &DMStore{
			Error: err,
		}
	}
	return &DMStore{
		BeSecretStore: dmSec,
		BeDataStore:   dmDb,
	}
}

// Initializing DMStore struct with object and global var
func InitDMStore(conf *utils.VapusAIStudioConfig) {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	if DMStoreManager == nil || DMStoreManager.BeSecretStore == nil || DMStoreManager.BeDataStore == nil {
		DMStoreManager = newDMStore(conf)
	}
}
