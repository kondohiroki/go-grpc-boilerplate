package middleware

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic along with request ID and other relevant details
				requestID := getRequestIDFromContext(ctx)
				logger.Error("Unexpected panic occurred in Unary call",
					zap.Any("panic", r),
					zap.String("x-request-id", requestID),
					zap.String("method", info.FullMethod),
				)
				err = status.Errorf(codes.Internal, "Unexpected error occurred")
			}
		}()
		return handler(ctx, req)
	}
}

func StreamRecoveryInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic along with request ID and other relevant details
				requestID := getRequestIDFromContext(stream.Context())
				logger.Error("Unexpected panic occurred in Stream call",
					zap.Any("panic", r),
					zap.String("x-request-id", requestID),
					zap.String("method", info.FullMethod),
				)
				err = status.Errorf(codes.Internal, "Unexpected error occurred")
			}
		}()
		return handler(srv, stream)
	}
}
