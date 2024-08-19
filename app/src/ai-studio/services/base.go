package services

import (
	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
)

type DMServices struct {
	DMStore *dmstores.DMStore
	*AIStudioServices
}

var DMServicesManager *DMServices
var helperLogger zerolog.Logger

func newDMServices(dmstore *dmstores.DMStore) *DMServices {
	return &DMServices{
		DMStore: dmstore,
	}
}

func InitDMServices(dmstore *dmstores.DMStore) {
	InitAIStudioServices(dmstore)
	helperLogger = pkgs.GetSubDMLogger(pkgs.SVCS, "helpers")
	if DMServicesManager == nil {
		DMServicesManager = newDMServices(dmstore)
		DMServicesManager.AIStudioServices = AIStudioServicesManager
	}
}
