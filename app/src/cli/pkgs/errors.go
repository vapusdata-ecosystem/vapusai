package pkg

import "errors"

var (
	ErrNoArgs      = errors.New("no arguments provided for the command, please provide the required arguments")
	ErrInvalidArgs = errors.New("invalid arguments provided for the command, please provide the required arguments")

	ErrLoadingGrpcCert         = errors.New("error loading grpc cert for grpc client, please check the cert path provided in cli config")
	ErrGrpcConn                = errors.New("error establishing grpc connection with the datamesh server, please check the server address and authtoken provided in cli config")
	ErrNoDataMeshFound         = errors.New("no datamesh found for the provided meshid")
	ErrNoDomainNodeFound       = errors.New("no domain node found for the provided domain node id")
	ErrNoDataNodeFound         = errors.New("no domain node found for the provided domain node id")
	ErrInvalidInterfaceGoal    = errors.New("invalid interface action provided, please provide a valid action")
	ErrNoCurrentContext        = errors.New("no current context found, please set a current context")
	ErrInvalidAction           = errors.New("invalid action provided, please provide a valid action")
	ErrMissingDataProductLogin = errors.New("no data product provided for login")
	ErrInvalidDataWorkerType   = errors.New("invalid data worker type provided, please provide a valid data worker type")
	ErrMetaData404             = errors.New("metadata not found for the provided data source")
)
