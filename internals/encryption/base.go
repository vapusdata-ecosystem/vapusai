package encrytion

import (
	"context"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

var encryptLogger zerolog.Logger

// TO:DO Make it more generic to handle different type of claims
type JwtAuthService interface {
	GenerateStudioJwt(claims *models.AIStudioSpaceClaims) (string, error)
	Parse(tokenString string) (*models.AIStudioSpaceClaims, error)
	Validate(tokenString string) (*jwt.Token, error)
	ValidateStudioAccessToken(tokenString string) (map[string]string, error)
}

type JWTAuthn struct {
	PublicJWTKey     string `validate:"required" yaml:"publicJwtKey" json:"publicJwtKey"`
	PrivateJWTKey    string `validate:"required" yaml:"privateJwtKey" json:"privateJwtKey"`
	SigningAlgorithm string `validate:"required" yaml:"signingAlgorithm" json:"signingAlgorithm"`
	JwtAudience      string `validate:"required" yaml:"jwtAudience" json:"jwtAudience"`
}

type jwtAuthOpts func(jo *JWTAuthn)

type VapusDataJwtAuthn struct {
	Opts *JWTAuthn
	JwtAuthService
}

var JwtTokenIssuer = "VapusData"
var JwtTokenAudience = "*.vapusdata.com"
var VapusPlatformTokenSubject = "VapusData access token"
var JwtDomainScope = "DomainScope"
var JwtPlatformScope = "platformScope"
var JwtDataProductScope = "dataProductScope"
var JwtCtxClaimKey = "vapusPlatformJwtClaim"
var JwtDPCtxClaimKey = "vapusPlatformJwtClaim"
var JwtClaimRoleSeparator = "|"

func SetCtxClaim(ctx context.Context, claim map[string]string) context.Context {
	return context.WithValue(ctx, JwtCtxClaimKey, claim)
}

func SetDataProductCtxClaim(ctx context.Context, claim map[string]string) context.Context {
	return context.WithValue(ctx, JwtDPCtxClaimKey, claim)
}

func GetCtxClaim(ctx context.Context) (map[string]string, bool) {
	val, ok := ctx.Value(JwtCtxClaimKey).(map[string]string)
	return val, ok
}

func GetDPtxClaim(ctx context.Context) (map[string]string, bool) {
	val, ok := ctx.Value(JwtDPCtxClaimKey).(map[string]string)
	return val, ok
}

func NewAuthzWithConfig(path string) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	jwtAuthnSecrets, err := LoadJwtAuthnSecrets(path)
	if err != nil {
		return nil, err
	}
	return NewAuthz(jwtAuthnSecrets)
}

func LoadJwtAuthnSecrets(path string) (*JWTAuthn, error) {
	cf, err := dmutils.ReadBasicConfig(dmutils.GetConfFileType(path), path, &JWTAuthn{})
	if err != nil {
		encryptLogger.Info().Msgf("Error loading jwt authn secrets: %v", err)
		return nil, err
	}
	return cf.(*JWTAuthn), err
}

func NewAuthz(opts *JWTAuthn) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	obj := &VapusDataJwtAuthn{
		Opts: opts,
	}
	switch opts.SigningAlgorithm {
	case mpb.EncryptionAlgo_ECDSA.String():
		val, err := NewECDSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val
		return obj, nil
	case mpb.EncryptionAlgo_RSA.String():
		val, err := NewRSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val
		return obj, nil
	default:
		return nil, ErrInvalidJWT
	}

}
