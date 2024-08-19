package pkgs

import (
	"github.com/vapusdata-ecosystem/vapusai-studio/aistudio/utils"
	authn "github.com/vapusdata-ecosystem/vapusai-studio/internals/authn"
)

type AuthnService struct {
	*authn.Authenticator
	Auth, Callback string
}

var AuthnManager *AuthnService
var authLogger = getSubPkgLogger(SVCS, "Authentication")

func InitAuthnManager(params *authn.AuthnSecrets) {
	if AuthnManager == nil {
		authnSrv, err := authn.NewAuthenticator(params.OIDCSecrets, params.AuthnMethod)
		if err != nil {
			authLogger.Panic().Err(err).Msgf(utils.ErrAuthenticatorInitFailed.Error())
		}
		AuthnManager = &AuthnService{
			Authenticator: authnSrv,
		}
	}
}
