package pkgs

import (
	grpctools "github.com/vapusdata-ecosystem/vapusai-studio/internals/grpctools"
)

var requestValidator *grpctools.StudioRequestValidator

func initRequestValidator() {
	validator, err := grpctools.NewGRPCValidator()
	if err != nil {
		pkgLogger.Panic().Err(err).Msg("Error while loading validator")
	}
	requestValidator = validator
}

func GetRequestValidator() *grpctools.StudioRequestValidator {
	if requestValidator == nil {
		initRequestValidator()
	}
	return requestValidator
}
