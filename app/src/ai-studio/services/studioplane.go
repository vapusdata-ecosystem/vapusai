package services

import dmstores "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/datastoreops"

type AIStudioServices struct {
	DMStore *dmstores.DMStore
}

var AIStudioServicesManager *AIStudioServices

func NewAIStudioServices(dmstore *dmstores.DMStore) *AIStudioServices {
	return &AIStudioServices{
		DMStore: dmstore,
	}
}

func InitAIStudioServices(dmstore *dmstores.DMStore) {
	if AIStudioServicesManager == nil {
		AIStudioServicesManager = NewAIStudioServices(dmstore)
	}
}
