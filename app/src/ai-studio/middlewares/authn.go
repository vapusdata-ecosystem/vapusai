package middlewares

import (
	"context"

	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusai-studio/internals/encryption"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Initiate authenticator function for DataMesh JWT Authenication
// This function will be used as a middleware to authenticate the request
func AuthnMiddleware(ctx context.Context) (context.Context, error) {
	methodName, _ := grpc.Method(ctx)
	if !needAuthn(methodName) {
		return ctx, nil
	}
	logger = pkgs.GetSubDMLogger("Middleware", "Authn")
	logger.Info().Msgf("Authenticating request for method - %v", methodName)
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		logger.Err(err).Msg("error while obtaining token from request header")
		return nil, status.Error(codes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	return AiStudioAuthn(ctx, token)
}

func AiStudioAuthn(ctx context.Context, token string) (context.Context, error) {
	parsedClaims, err := pkgs.VapusAuth.ValidateStudioAccessToken(token)
	if err != nil {
		logger.Err(err).Msg("error while validating data product access access token from request header")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	// logger.Info().Msgf("parsed domain claims - %v", parsedClaims)
	// _, role, err := dmstores.DMStoreManager.GetDomainUser(ctx, parsedClaims[encryption.ClaimUserIdKey], parsedClaims[encryption.ClaimDomainKey])
	// logger.Info().Msgf("role - %v", role)
	// logger.Info().Msgf("parsedClaims - %v", parsedClaims)
	// if err != nil {
	// 	logger.Err(err).Msgf("error while validating access token against the claim domain, user - %v, domain - %v", parsedClaims[encryption.ClaimUserIdKey], parsedClaims[encryption.ClaimDomainKey])
	// 	return nil, status.Error(codes.Unauthenticated, utils.ErrDomain404.Error())
	// } else if role != parsedClaims[encryption.ClaimDomainRoledKey] {
	// 	logger.Err(err).Msgf("error while validating access token against the claim domain, user - %v, domain - %v, role - %v", parsedClaims[encryption.ClaimUserIdKey], parsedClaims[encryption.ClaimDomainKey], parsedClaims[encryption.ClaimDomainRoledKey])
	// 	return nil, status.Error(codes.Unauthenticated, utils.ErrUserDomain404.Error())
	// }
	return encryption.SetCtxClaim(ctx, parsedClaims), nil
}

func needAuthn(funcName string) bool {
	return true
}
