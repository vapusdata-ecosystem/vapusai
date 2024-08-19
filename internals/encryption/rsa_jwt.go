package encrytion

import (
	"crypto/rsa"

	jwt "github.com/golang-jwt/jwt/v5"
	dmerrors "github.com/vapusdata-ecosystem/vapusai-studio/internals/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

type RSAJwt interface {
	GenerateStudioJwt(claims *models.AIStudioSpaceClaims) (string, error)
	ParseStudioJwtToken(tokenString string) (*models.AIStudioSpaceClaims, error)
	ValidateStudioJwtToken(tokenString string) (*jwt.Token, error)
	ValidateStudioAccessToken(tokenString string) (map[string]string, error)
}

type RSAManager struct {
	opts        *JWTAuthn
	ParsedPvKey *rsa.PrivateKey
	ParsedPbKey *rsa.PublicKey
}

var rsaSigningAlgoMap = map[string]*jwt.SigningMethodRSA{
	"P-521": jwt.SigningMethodRS512,
	"P-384": jwt.SigningMethodRS384,
	"P-256": jwt.SigningMethodRS256,
}

// NewRSAJwtAuthn creates a new RSA JWT Authn object with the given options.
// It returns the RSAJwt interface. It logs an error if the private key is not parsed.
func NewRSAJwtAuthn(opts *JWTAuthn) (*RSAManager, error) {
	parsedPvKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(opts.PrivateJWTKey))
	if err != nil || parsedPvKey == nil {
		dmlogger.CoreLogger.Err(err).Msg("Error parsing RSA private key")
		return nil, err
	}

	// TODO: Add validation for public key and private key after parsing
	return &RSAManager{
		opts:        opts,
		ParsedPvKey: parsedPvKey,
		ParsedPbKey: &parsedPvKey.PublicKey,
	}, nil
}

func (e *RSAManager) GenerateStudioJwt(claims *models.AIStudioSpaceClaims) (string, error) {
	dmlogger.CoreLogger.Info().Msgf("Generating JWT token for claim %v", claims)
	token := jwt.NewWithClaims(rsaSigningAlgoMap[e.ParsedPvKey.N.String()], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return dmutils.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *RSAManager) Parse(tokenString string) (*models.AIStudioSpaceClaims, error) {
	token, err := e.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.AIStudioSpaceClaims); !ok {
		dmlogger.CoreLogger.Err(ErrInvalidJWTClaims).Msg("Invalid JWT claims")
		return nil, dmerrors.DMError(ErrInvalidJWTClaims, nil)
	} else {
		return claims, nil

	}
}

func (e *RSAManager) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AIStudioSpaceClaims{}, func(token *jwt.Token) (interface{}, error) {
		return e.ParsedPbKey, nil
	})
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msg(ErrParsingJWT.Error())
		return nil, dmerrors.DMError(ErrParsingJWT, err)
	} else if !token.Valid {
		dmlogger.CoreLogger.Err(err).Msg("Invalid JWT token")
		return nil, dmerrors.DMError(ErrInvalidJWT, nil)
	} else {
		return token, nil
	}
}

func (e *RSAManager) ValidateStudioAccessToken(tokenString string) (map[string]string, error) {
	claim, err := e.Parse(tokenString)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error while parsing and validating auth token")
		return nil, err
	}
	resp := FlatJwtClaims(claim, "||")
	if resp == nil {
		dmlogger.CoreLogger.Error().Msgf("invalid Claim parsed from the token")
		return nil, dmerrors.DMError(ErrInvalidUserAuthentication, nil)
	}
	dmlogger.CoreLogger.Info().Msgf("flatClaims - %v", resp)
	return resp, nil
}
