package pkg

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetBearerCtx(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}
