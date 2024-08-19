package cmd

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
)

type PlatformInstanceClients struct {
	*plclient.VapusPlatformClient
}

var (

	// RootCmd is the root command for vapusdt
	rootCmd                                                       *cobra.Command
	cfgFile                                                       string
	debugLogFlag                                                  bool
	vapusGlobals                                                  *GlobalsPersists
	logger                                                        zerolog.Logger
	currentIdToken, currentAccessToken, currentProductAccessToken string = "currentIdToken", "currentAccessToken", "currentProductAccessToken"
	GlobalVar                                                     string
	action, file, search                                          string
	la                                                            bool
)

var ignoreConnMap = map[string]bool{
	configCmd:  true,
	clearCmd:   true,
	explainCmd: true,
}

type GlobalsPersists struct {
	CurrentContext    string
	logger            zerolog.Logger
	cfgFile           string
	cfgDir            string
	debugLogFlag      bool
	AgentsActions     map[string][]interface{}
	AgentsUtilities   map[string][]interface{}
	AgentsReflexes    map[string][]interface{}
	AgentInterfaceMap map[string]string

	*plclient.VapusPlatformClient
	currentIdToken, currentAccessToken string
	ctx                                context.Context
}

// commands and resource var/constants
const (
	contextsCmd         = "context"
	configCmd           = "config"
	datameshResource    = "datamesh"
	dataSourceResource  = "datasources"
	domainResource      = "domains"
	dataProductResource = "dataproducts"
	dataWorkerResource  = "dataworkers"
	dataCatalogResource = "datacatalogs"
	explainCmd          = "explain"
	operationcmd        = "operations"
	authAction          = "auth"
	genTemplate         = "gen-template"
	domainAuthCmd       = "domain"
	dataProductAuthCmd  = "dataproduct"
	clearCmd            = "clear"
	initializeCmd       = "init"
	interfaceCmd        = "interface"
	connectCmd          = "connect"
	SpecsCmd            = "spec"
	generateCmd         = "generate"
	loginCmd            = "login"
	accountResource     = "account"
	userResource        = "users"
)
