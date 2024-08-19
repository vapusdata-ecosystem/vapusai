package encrytion

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	jwt "github.com/golang-jwt/jwt/v5"
	dmerrors "github.com/vapusdata-ecosystem/vapusai-studio/internals/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

var DefaultECDSElliptic string = "P-521"

type ECDSAKeys struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
	EllipticCurve elliptic.Curve
}

type ECDSAManager struct {
	opts        *JWTAuthn
	ParsedPvKey *ecdsa.PrivateKey
	ParsedPbKey *ecdsa.PublicKey
}

var ecdsaSigningAlgoMap = map[string]*jwt.SigningMethodECDSA{
	"P-521": jwt.SigningMethodES512,
	"P-384": jwt.SigningMethodES384,
	"P-256": jwt.SigningMethodES256,
}

var ellipticCurveMap = map[string]elliptic.Curve{
	"P-256": elliptic.P256(),
	"P-384": elliptic.P384(),
	"P-521": elliptic.P521(),
}

func GenerateECDSAKeys(curve string) (*ECDSAKeys, error) {
	eCurve := ellipticCurveMap[curve]
	privKey, err := ecdsa.GenerateKey(eCurve, rand.Reader)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error generating ECDSA private key with elliptic curve %v", curve)
		return nil, err
	}
	return &ECDSAKeys{
		PrivateKey:    privKey,
		PublicKey:     &privKey.PublicKey,
		EllipticCurve: eCurve,
	}, nil
}

// NewECDSAJwtAuthn creates a new ECDSA JWT Authn object with the given options.
// It returns the ECDSAJwt interface. It logs an error if the private key is not parsed.
func NewECDSAJwtAuthn(opts *JWTAuthn) (*ECDSAManager, error) {
	parsedPvKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(opts.PrivateJWTKey))
	if err != nil || parsedPvKey == nil {
		dmlogger.CoreLogger.Err(err).Msg("Error parsing ECDSA private key")
		return nil, err
	}

	// TODO: Add validation for public key and private key after parsing
	return &ECDSAManager{
		opts:        opts,
		ParsedPvKey: parsedPvKey,
		ParsedPbKey: &parsedPvKey.PublicKey,
	}, nil
}

func (e *ECDSAManager) GenerateStudioJwt(claims *models.AIStudioSpaceClaims) (string, error) {
	dmlogger.CoreLogger.Info().Msgf("Generating JWT token for claim %v", claims)
	token := jwt.NewWithClaims(ecdsaSigningAlgoMap[e.ParsedPvKey.Curve.Params().Name], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return dmutils.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *ECDSAManager) Parse(tokenString string) (*models.AIStudioSpaceClaims, error) {
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

func (e *ECDSAManager) Validate(tokenString string) (*jwt.Token, error) {
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

func (e *ECDSAManager) ValidateStudioAccessToken(tokenString string) (map[string]string, error) {
	claim, err := e.Parse(tokenString)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error while parsing and validating auth token")
		return nil, err
	}
	dmlogger.CoreLogger.Info().Msgf("parsed domain claims - %v", claim)
	resp := FlatJwtClaims(claim, "||")
	if resp == nil {
		dmlogger.CoreLogger.Error().Msgf("invalid Claim parsed from the token")
		return nil, dmerrors.DMError(ErrInvalidUserAuthentication, nil)
	}
	return resp, nil
}
