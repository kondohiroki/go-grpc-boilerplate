package server

import (
	"fmt"

	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/middleware"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UserServiceServer

	// Add more services here
}

func (s *Server) registerWithServer(sv *grpc.Server) {
	pb.RegisterUserServiceServer(sv, s)

	// Register more services here
}

// initOptions initializes the gRPC server options
func initOptions() []grpc.ServerOption {
	maxSendSize := config.GetConfig().GRPCServer.MaxSendMsgSize * 1024 * 1024 // Convert MB to bytes
	maxRecvSize := config.GetConfig().GRPCServer.MaxRecvMsgSize * 1024 * 1024 // Convert MB to bytes

	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(maxSendSize),
		grpc.MaxRecvMsgSize(maxRecvSize),

		// Register interceptors here
		grpc.ChainUnaryInterceptor(
			middleware.UnaryRequestIDInterceptor,
			middleware.UnaryLoggingInterceptor(logger.Log),
			middleware.UnaryRecoveryInterceptor(logger.Log),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamRequestIDInterceptor,
			middleware.StreamLoggingInterceptor(logger.Log),
			middleware.StreamRecoveryInterceptor(logger.Log),
		),
	}

	if config.GetConfig().GRPCServer.UseTLS {
		logger.Log.Info("Using TLS for gRPC server")

		// Load certificates
		certFile := config.GetConfig().GRPCServer.TLSCertFile
		keyFile := config.GetConfig().GRPCServer.TLSKeyFile
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

		if err != nil {
			logger.Log.Fatal(fmt.Sprintf("Failed to load TLS: %v", err))
		}
		opts = append(opts, grpc.Creds(creds))
	}

	return opts
}

func NewGRPCServer() (*grpc.Server, error) {
	// Register interceptors here
	opts := initOptions()

	// Create a new gRPC server with the interceptors
	s := grpc.NewServer(opts...)

	// Initialize and register the combined gRPC server struct
	serverInstance := &Server{}
	serverInstance.registerWithServer(s)

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus(config.GetConfig().App.NameSlug, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// Register the reflection service on gRPC server.
	if config.GetConfig().GRPCServer.UseReflection {
		logger.Log.Info("Register reflection service")
		reflection.Register(s)
	}
	return s, nil
}
