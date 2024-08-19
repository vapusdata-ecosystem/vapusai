package models

import (
	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
)

type ConnectionOpts struct {
	DataSourceEngine      string                 `validate:"required" json:"dataSourceEngine" yaml:"dataSourceEngine" toml:"dataSourceEngine"`
	Params                map[string]interface{} `json:"params" yaml:"params" toml:"params"`
	DataSourceSvcProvider string                 `json:"dataSourceSvcProvider" yaml:"dataSourceSvcProvider" toml:"dataSourceSvcProvider"`
	Credentials           *mpb.VapusCredentials  `validate:"required" json:"credentials" yaml:"credentials" toml:"credentials"`
	URL                   string                 `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Port                  int64                  `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
}
