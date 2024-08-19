package authn

import (
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

type AuthnSecrets struct {
	AuthnMethod string       `json:"authnMethod" yaml:"authnMethod"`
	OIDCSecrets *OIDCSecrets `json:"oidcSecrets" yaml:"oidcSecrets"`
}

type OIDCSecrets struct {
	Domain       string `json:"domain" yaml:"domain"`
	ClientID     string `json:"clientId" yaml:"clientId"`
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
	LoginURL     string `json:"login" yaml:"login"`
	CallbackURI  string `json:"callbackUri" yaml:"callbackUri"`
}

// GetOIDCSecretStore function to read the OIDC secrets from the file provided.
func LoadAuthnSecrets(path string) (*AuthnSecrets, error) {
	// Read the file containing the OIDC secrets
	cf, err := dmutils.ReadBasicConfig(dmutils.GetConfFileType(path), path, &AuthnSecrets{})

	return cf.(*AuthnSecrets), err
}
