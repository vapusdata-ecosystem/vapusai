package dmcontrollers

type VapusAIStudioController struct {
	*AIStudioPlane
}

var VapusAIStudioManager *VapusAIStudioController

func NewVapusAIStudioController() *VapusAIStudioController {
	InitVapusAIStudioController()
	return &VapusAIStudioController{
		AIStudioPlane: AIStudioPlaneManager,
	}
}

func InitVapusAIStudioController() {
	if AIStudioPlaneManager == nil {
		AIStudioPlaneManager = NewAIStudioPlane()
	}
}
