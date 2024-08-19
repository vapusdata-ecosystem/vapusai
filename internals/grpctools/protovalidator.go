package pbtools

import (
	protovalidate "github.com/bufbuild/protovalidate-go"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// DMValidator struct to hold the validator
type StudioRequestValidator struct {
	Validator *protovalidate.Validator
}

// NewDMValidator function to create a new validator
func NewGRPCValidator() (*StudioRequestValidator, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	return &StudioRequestValidator{
		Validator: v,
	}, nil
}

// Validate function to validate the data
func (v *StudioRequestValidator) Validate(data interface{}) error {

	if err := v.Validator.Validate(data.(protoreflect.ProtoMessage)); err != nil {
		return err
	}
	return nil
}
