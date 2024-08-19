package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AIStudioSpaceClaims struct {
	jwt.RegisteredClaims
	Scope *StudioScope `validate:"required" json:"scope"`
}

type StudioScope struct {
	Space  string `validate:"required" json:"space"`
	Role   string `validate:"required" json:"role"`
	UserId string `validate:"required" json:"userId"`
}

func (x *StudioScope) Validate() error {
	validator := validator.New()
	return validator.Struct(x)
}
