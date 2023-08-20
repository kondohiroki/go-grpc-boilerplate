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
	// Add more gRPC services here
}

func NewgRPCServer() (*grpc.Server, error) {
	// Register interceptors here
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	}

	// Create a new gRPC server with the interceptors
	s := grpc.NewServer(opts...)

	// Register health check service
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	// Register the gRPC server and more gRPC services here
	pb.RegisterUsersServer(s, &Server{})

	// Register the reflection service on gRPC server.
	reflection.Register(s)

	return s, nil
}

// Custom recovery function to handle panics
func grpcPanicRecoveryHandler(p interface{}) (err error) {
	logger.Log.Error("Unexpected panic occurred", zap.Any("panic", p))
	return status.Errorf(codes.Internal, "Unexpected error occurred")
}
