package server

import (
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/middleware"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
	"google.golang.org/grpc"
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

// Register interceptors (i.e. middleware) here
func initOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.UnaryRequestIDInterceptor,
			middleware.UnaryLoggingInterceptor(logger.Log),
			middleware.UnaryRecoveryInterceptor(logger.Log),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamLoggingInterceptor(logger.Log),
			middleware.StreamRecoveryInterceptor(logger.Log),
		),
	}
}

func NewGRPCServer() (*grpc.Server, error) {
	// Register interceptors here
	opts := initOptions()

	// Create a new gRPC server with the interceptors
	s := grpc.NewServer(opts...)

	// Register health check service
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	// Initialize and register the combined gRPC server struct
	serverInstance := &Server{}
	serverInstance.registerWithServer(s)

	// Register the reflection service on gRPC server.
	reflection.Register(s)

	return s, nil
}
