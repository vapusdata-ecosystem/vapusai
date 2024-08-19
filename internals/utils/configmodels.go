package utils

import (
	"log"
	filepath "path/filepath"
)

type VapusAIStudioConfig struct {
	Path                 string
	VapusBESecretStorage struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"vapusBECacheStorage"`
	ServerConfig struct {
		Port        int32  `yaml:"port"`
		ServiceName string `yaml:"serviceName"`
		Scheme      string `yaml:"scheme"`
		ExternalURL string `yaml:"external"`
	} `yaml:"serverConfig"`
	JWTAuthnSecrets struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"JWTAuthnSecrets"`
	AuthnSecrets struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"authnSecrets"`
	ServerCerts struct {
		MtlsEnabled     bool   `yaml:"mtlsEnabled"`
		PlainTlsEnabled bool   `yaml:"plainTlsEnabled"`
		Insecure        bool   `yaml:"insecure"`
		CaCertFile      string `yaml:"caCertFile"`
		ServerCertFile  string `yaml:"serverCertFile"`
		ServerKeyFile   string `yaml:"serverKeyFile"`
		ClientCertFile  string `yaml:"serverCertFile"`
		ClientKeyFile   string `yaml:"serverKeyFile"`
	} `yaml:"serverCerts"`
	AuthnMethod string `yaml:"authnMethod"`
}

func (sc *VapusAIStudioConfig) GetSecretStoragePath() string {
	log.Println("Secret storage path: ", filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath))
	log.Println("Secret storage path: ", sc.VapusBESecretStorage.FilePath)
	log.Println("Secret storage path: ", sc.Path)
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *VapusAIStudioConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.FilePath
}

func (sc *VapusAIStudioConfig) GetCachStoragePath() string {
	return sc.VapusBECacheStorage.FilePath
}

func (sc *VapusAIStudioConfig) GetJwtAuthSecretPath() string {
	return sc.JWTAuthnSecrets.FilePath
}

func (sc *VapusAIStudioConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAIStudioConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAIStudioConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *VapusAIStudioConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *VapusAIStudioConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}

func (sc *VapusAIStudioConfig) GetAuthnSecrets() string {
	return sc.AuthnSecrets.FilePath
}
