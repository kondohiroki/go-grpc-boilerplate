package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Get requestID from context
		requestID := getRequestIDFromContext(ctx)

		// Proceed with the original request
		resp, err := handler(ctx, req)

		// Log details
		logger.Info(info.FullMethod,
			zap.String("method", info.FullMethod),
			zap.String("x-request-id", requestID),
			zap.Duration("duration", time.Since(startTime)),
			zap.Any("request", req),
			zap.Any("response", resp),
			// Any other metadata you'd like to log
		)

		return resp, err
	}
}

func StreamLoggingInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		// This will not capture the actual request & response objects in streaming, as they're handled inside the handler.
		// To log those, you'd need more complex code to wrap the stream and capture each message.
		// Here we're just logging the method, id, and timings.

		// Get requestID from wrapped context of stream
		requestID := getRequestIDFromContext(stream.Context())

		// Proceed with the original request
		err := handler(srv, stream)

		// Log details
		logger.Info(info.FullMethod,
			zap.String("method", info.FullMethod),
			zap.String("x-request-id", requestID),
			zap.Bool("client_streaming", info.IsClientStream),
			zap.Bool("server_streaming", info.IsServerStream),
			zap.Duration("duration", time.Since(startTime)),
			// Any other metadata you'd like to log
		)

		return err
	}
}
