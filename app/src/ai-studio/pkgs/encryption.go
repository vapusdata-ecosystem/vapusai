package pkgs

import (
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	encryption "github.com/vapusdata-ecosystem/vapusai-studio/internals/encryption"
	models "github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
)

var VapusAuth *encryption.VapusDataJwtAuthn

func NewVapusAiStudioAuthz(params *encryption.JWTAuthn, validator *validator.Validate) (*encryption.VapusDataJwtAuthn, error) {
	obj, err := encryption.NewAuthz(params)
	if err != nil {
		pkgLogger.Err(err).Msg("Error while loading jwt authnetication config")
		return nil, err
	}
	if err := validator.Struct(obj.Opts); err != nil {
		pkgLogger.Err(err).Msg("Error while validating jwt config")
		return nil, err
	}

	return obj, nil
}

func InitAuthzManager(params *encryption.JWTAuthn, validator *validator.Validate) {
	var err error
	if VapusAuth == nil {
		VapusAuth, err = NewVapusAiStudioAuthz(params, validator)
		if err != nil {
			pkgLogger.Err(err).Msg("error initializing authn")
			panic(err)
		}
	}
}

func BuildVDPAClaim(userID, space, role string, validTill time.Time) (*models.AIStudioSpaceClaims, error) {

	return &models.AIStudioSpaceClaims{
		Scope: &models.StudioScope{
			UserId: userID,
			Space:  space,
			Role:   role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   encryption.VapusPlatformTokenSubject,
			Audience:  []string{encryption.JwtTokenAudience},
			ExpiresAt: jwt.NewNumericDate(validTill), // configurable tokens
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 2)),
		},
	}, nil
}
