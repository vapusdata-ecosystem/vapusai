package gcp

import (
	"fmt"

	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
)

type GcpConfig struct {
	ProjectID, Region, Zone string
	ServiceAccountKey       []byte
}

var (
	GCP_RES_TEMP = "projects/%s"
	GCP_SM_RES   = "secrets/%s"
)

var logger = dmlogger.CoreLogger

func (x *GcpConfig) GetGcpResource(location, resourceName, resourceValue string) string {
	return fmt.Sprintf("projects/%s/locations/%s/%s/%s", x.ProjectID, location, resourceName, resourceValue)
}
