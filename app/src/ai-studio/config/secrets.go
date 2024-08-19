package config

import (
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/authn"
	encrytion "github.com/vapusdata-ecosystem/vapusai-studio/internals/encryption"
	models "github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
)

type VapusArtifactStorage struct {
	Spec *models.ConnectionOpts `yaml:"spec"`
}

var VapusArtifactStorageManager *VapusArtifactStorage

var AuthnParams *authn.AuthnSecrets

var JwtAuthnParams *encrytion.JWTAuthn
