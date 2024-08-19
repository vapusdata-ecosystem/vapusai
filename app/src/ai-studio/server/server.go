package server

import (
	"context"
	"flag"

	"github.com/rs/zerolog"
	config "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/config"
	controllers "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/controllers"
	middlewares "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/middlewares"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/aistudio/v1alpha1"
	grpctools "github.com/vapusdata-ecosystem/vapusai-studio/internals/grpctools"

	interceptors "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	selector "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc "google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	debugLogFlag bool
	flagconfPath string
	configName   = "config/platform-service-config.yaml"
	bootConfig   = "config/vapusplatform-boot-config.yaml"
	logger       zerolog.Logger
)

func init() {
	flag.StringVar(&flagconfPath, "conf", "/data/vapusdata/platform/config", "config path, eg: --conf=/data/domain")
	flag.BoolVar(&debugLogFlag, "debug", false, "debug loggin, set it to true to enable the debug logs")
	flag.Parse()
	logger.Info().Msgf("Config root Path: %s", flagconfPath)
	packagesInit()
}

func initServer(grpcServer *grpctools.GRPCServer) {

	// Setup auth matcher.
	allButTheez := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	// Add unary and stream interceptors for prometheus
	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			selector.UnaryServerInterceptor(rpcauth.UnaryServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
			selector.UnaryServerInterceptor(middlewares.UnaryRequestValidator(), selector.MatchFunc(allButTheez)),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			selector.StreamServerInterceptor(rpcauth.StreamServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
		),
	}

	// Create a new GRPC server
	//First step is to configure the vapusPlatform server
	grpcServer.ServerPort = config.ServiceConfigManager.ServerConfig.Port

	// Initialize the server
	logger.Info().Msg("Configuring VapusData Platform Server")

	// Initialize the grpc server ops and net listner for the server
	grpcServer.ConfigureGrpcServer(opts, debugLogFlag)
	logger.Info().Msg("VapusData Platform Server configured successfully.")
	if grpcServer.GrpcServ == nil {
		logger.Info().Msg("Failed to initialize VapusData Platform Server")
	}
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer.GrpcServ, healthcheck)
	// Register the VapusData Platform and Node controller
	reflection.Register(grpcServer.GrpcServ)
	pb.RegisterVapusAiStudioServer(grpcServer.GrpcServ, controllers.NewAIStudioPlane())

	logger.Info().Msgf("Server configured at - %v", grpcServer.ServerPort)

}

func GrpcServer() *grpctools.GRPCServer {
	var grpcServer *grpctools.GRPCServer
	var serverOpts []grpctools.GrpcOptions
	// Initialize the service configuration
	if config.ServiceConfigManager.ServerCerts.MtlsEnabled {
		logger.Info().Msg("Configuring VapusData Platform Server with MTLS connection")
		serverOpts = append(serverOpts, grpctools.WithMtls(config.ServiceConfigManager.GetMtlsCerts()))
	} else if config.ServiceConfigManager.ServerCerts.PlainTlsEnabled {
		logger.Info().Msg("Configuring VapusData Platform Server with PlainTLS connection")
		serverOpts = append(serverOpts, grpctools.WithPlainTls(config.ServiceConfigManager.GetPlainTlsCerts()))
	} else {
		logger.Info().Msg("Configuring VapusData Platform Server with insecure connection")
		serverOpts = append(serverOpts, grpctools.WithInsecure(true))
	}
	grpcServer = grpctools.NewGRPCServer(serverOpts...)
	initServer(grpcServer)
	return grpcServer
}
