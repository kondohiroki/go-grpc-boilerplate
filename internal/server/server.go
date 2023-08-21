package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UsersServer

	// Add more services here
}

func (s *Server) registerWithServer(sv *grpc.Server) {
	pb.RegisterUsersServer(sv, s)

	// Register more services here
}

// Register interceptors (i.e. middleware) here
func initOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
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

// Custom recovery function to handle panics
func grpcPanicRecoveryHandler(p interface{}) (err error) {
	logger.Log.Error("Unexpected panic occurred", zap.Any("panic", p))
	return status.Errorf(codes.Internal, "Unexpected error occurred")
}
