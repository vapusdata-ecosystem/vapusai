package pbtools

import (
	"context"
	"fmt"

	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
	dmerrors "github.com/vapusdata-ecosystem/vapusai-studio/internals/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

var ErrorFinal = "errorFinal"

// HandleGrpcError is a utility function to handle grpc errors
func HandleGrpcError(err error, code grpccodes.Code) error {
	e := err.(dmerrors.Error)
	return grpcstatus.Error(code, e.Error())
}

func HandleCtxCustomMessage(ctx context.Context, msgType string, msg ...string) context.Context {
	if ctx.Value(dmutils.CUSTOM_MESSAGE) == nil {
		cm := map[string]interface{}{msgType: msg}
		return context.WithValue(ctx, dmutils.CUSTOM_MESSAGE, cm)
	}
	cm := ctx.Value(dmutils.CUSTOM_MESSAGE).(map[string]interface{})
	cm[msgType] = append(cm[msgType].([]string), msg...)
	return context.WithValue(ctx, dmutils.CUSTOM_MESSAGE, cm)
}

// HandleResponse is a utility function to handle the base response
func HandleDMResponse(ctx context.Context, opts ...string) *mpb.DMResponse {
	if ctx.Value(dmutils.CUSTOM_MESSAGE) == nil {
		return nil
	}
	cm := ctx.Value(dmutils.CUSTOM_MESSAGE).(map[string]interface{})
	return &mpb.DMResponse{
		Message:      opts[0],
		DmStatusCode: opts[1],
		CustomMessage: func(cms map[string]interface{}) []*mpb.MapList {
			var cm []*mpb.MapList
			for k, v := range cms {
				cm = append(cm, &mpb.MapList{
					Key:    k,
					Values: v.([]string),
				})
			}
			return cm
		}(cm),
	}
}

func GetSvcDns(svcName string, namespace string, port int64) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", svcName, namespace, port)
}
