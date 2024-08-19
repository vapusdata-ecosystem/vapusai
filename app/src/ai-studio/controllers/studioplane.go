package dmcontrollers

import (
	"github.com/rs/zerolog"

	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
	dmsvc "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/services"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/aistudio/v1alpha1"
)

type AIStudioPlane struct {
	pb.UnimplementedVapusAiStudioServer
	DMServices *dmsvc.DMServices
	Logger     zerolog.Logger
}

var AIStudioPlaneManager *AIStudioPlane

func NewAIStudioPlane() *AIStudioPlane {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIStudioPlane")

	l.Info().Msg("AIStudioPlane Controller initialized")
	return &AIStudioPlane{
		Logger:     l,
		DMServices: dmsvc.DMServicesManager,
	}
}

func InitAIStudioPlane() {
	if AIStudioPlaneManager == nil {
		AIStudioPlaneManager = NewAIStudioPlane()
	}
}
