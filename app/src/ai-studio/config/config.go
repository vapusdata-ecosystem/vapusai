package config

import (
	validator "github.com/go-playground/validator/v10"
	pkgs "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
)

var ServiceConfigManager *utils.VapusAIStudioConfig
var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func newServiceConfig(configRoot, path string) *utils.VapusAIStudioConfig {
	return LoadServiceConfig(configRoot, path)
}

func InitServiceConfig(configRoot, path string) {
	if ServiceConfigManager == nil {
		ServiceConfigManager = newServiceConfig(configRoot, path)
	}
}

func LoadServiceConfig(configRoot, path string) *utils.VapusAIStudioConfig {
	// Read the service configuration from the file
	pkgs.DmLogger.Info().Msgf("Reading service configuration with path - %v ", path)

	cf, err := utils.ReadBasicConfig(utils.GetConfFileType(path), path, &utils.VapusAIStudioConfig{})
	if err != nil {
		pkgs.DmLogger.Panic().Err(err).Msg("error while loading service config")
		return nil
	}

	svcConf := cf.(*utils.VapusAIStudioConfig)
	svcConf.Path = configRoot
	return svcConf
}
