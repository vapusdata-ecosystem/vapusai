package pbtools

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	prometheus "github.com/prometheus/client_golang/prometheus"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/internals/utils"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	grpc_insecure "google.golang.org/grpc/credentials/insecure"
)

// GRPCServer struct with configuration options to start and run the grpc server
type GRPCServer struct {
	ServerPort    int32
	ServerAddress string
	GrpcServ      *grpc.Server
	Netlis        net.Listener
	Logger        zerolog.Logger
	PlainTls      bool
	MTls          bool
	tlsCertFile   string
	tlsKeyFile    string
	tlsCaCertFile string
	Insecure      bool
}

type GrpcOptions func(*GRPCServer)

// NewGRPCServer creates a new grpc server with configuration options
func NewGRPCServer(opts ...GrpcOptions) *GRPCServer {
	grpcServer := &GRPCServer{}
	for _, o := range opts {
		o(grpcServer)
	}
	return grpcServer
}

func WithMtls(caFile, serverFile, keyFile string) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.tlsCaCertFile = caFile
		opt.tlsCertFile = serverFile
		opt.tlsKeyFile = keyFile
		opt.MTls = true
	}
}

func WithPlainTls(serverFile, keyFile string) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.tlsCertFile = serverFile
		opt.tlsKeyFile = keyFile
		opt.PlainTls = true
	}
}

func WithInsecure(insecure bool) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.Insecure = insecure
	}
}

// getServerAddress returns the server address
func (gs *GRPCServer) setServerAddress() {
	// if server address is not provided, use localhost
	if gs.ServerAddress == dmutils.EMPTYSTR {
		gs.ServerAddress = fmt.Sprintf(":%d", gs.ServerPort)
	}
}

// NewGrpcServer creates a new grpc server with confirguration and other options
// This function will return the grpc server and net listner
func (gs *GRPCServer) ConfigureGrpcServer(opts []grpc.ServerOption, debugFlag bool) {
	var err error
	// Set the server address
	gs.setServerAddress()

	// If tls is enabled, generate the credentials
	opts = gs.initServerCreds(opts)

	grpcServer := grpc.NewServer(opts...)

	//register the grpc server with prometheus
	grpc_prometheus.EnableHandlingTimeHistogram(grpc_prometheus.WithHistogramBuckets(prometheus.ExponentialBucketsRange(0.5, 70, 9)))
	grpc_prometheus.Register(grpcServer)
	gs.Netlis, err = net.Listen("tcp", gs.ServerAddress)
	if err != nil {
		gs.Logger.Fatal().Err(err)
	}
	gs.GrpcServ = grpcServer
}

func (gs *GRPCServer) Run() {

	err := gs.GrpcServ.Serve(gs.Netlis)

	if err != nil {
		gs.Logger.Panic().Err(err)
	}
}

func (gs *GRPCServer) initServerCreds(opts []grpc.ServerOption) []grpc.ServerOption {

	// If Plain tls is enabled, generate the credentials
	if gs.PlainTls {
		return gs.initPlainTls(opts)
	} else if gs.MTls {
		return gs.initMtls(opts)
	} else if gs.Insecure {
		return gs.initInsecure(opts)
	}
	return opts
}

// initPlainTls initializes the Plain TLS credentials for the GRPCServer.
// It generates the credentials using the provided TLS certificate and key files.
// If the certificate or key file is not provided, it falls back to the default files.
// It returns the updated list of server options with the Plain TLS credentials.
func (gs *GRPCServer) initPlainTls(opts []grpc.ServerOption) []grpc.ServerOption {
	gs.Logger.Info().Msg("Generating Plain TLS credentials")
	if gs.tlsCertFile == "" {
		gs.Logger.Info().Msg("Using default server cert file")
		gs.tlsCertFile = dmutils.DEFAULTSERVERCERTTLS
	}
	if gs.tlsKeyFile == "" {
		gs.Logger.Info().Msg("Using default server key file")
		gs.tlsKeyFile = dmutils.DEFAULTSERVERKEYTLS
	}
	creds, err := credentials.NewServerTLSFromFile(gs.tlsCertFile, gs.tlsKeyFile)
	if err != nil {
		gs.Logger.Fatal().Err(err).Msg("Failed to generate credentials")
	}
	return append(opts, grpc.Creds(creds))
}

// initMtls initializes the MTLS (Mutual TLS) credentials for the GRPCServer.
// It generates the MTLS credentials using the provided CA certificate, server certificate, and key.
// The function returns a list of GRPC server options with the MTLS credentials appended.
func (gs *GRPCServer) initMtls(opts []grpc.ServerOption) []grpc.ServerOption {
	gs.Logger.Info().Msg("Generating MTLS credentials")
	caPem, err := os.ReadFile(gs.tlsCaCertFile)
	if err != nil {
		gs.Logger.Fatal().Err(err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		gs.Logger.Fatal().Err(err)
	}

	// read server cert & key
	serverCert, err := tls.LoadX509KeyPair(gs.tlsCertFile, gs.tlsKeyFile)
	if err != nil {
		gs.Logger.Fatal().Err(err)
	}

	// configuration of the certificate what we want to
	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	//create tls certificate
	tlsCredentials := credentials.NewTLS(conf)
	return append(opts, grpc.Creds(tlsCredentials))
}

func (gs *GRPCServer) initInsecure(opts []grpc.ServerOption) []grpc.ServerOption {
	return append(opts, grpc.Creds(grpc_insecure.NewCredentials()))
}
