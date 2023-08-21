package middleware

import (
	"context"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Get requestID from context
		requestID := getRequestIDFromContext(ctx)

		// Get client IP
		clientIP := getClientIP(ctx)

		// Extract method and service name
		methodName := getMethodName(info.FullMethod)
		serviceName := getServiceName(info.FullMethod)

		// Proceed with the original request
		resp, err := handler(ctx, req)

		// Log details
		logger.Info("unary request",
			zap.String("protocol", "grpc"),
			zap.String("service", serviceName),
			zap.String("method", methodName),
			zap.String("method_type", "unary"),
			zap.String("ip", clientIP),
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

		// Get requestID from wrapped context of stream
		requestID := getRequestIDFromContext(stream.Context())

		// Get client IP
		clientIP := getClientIP(stream.Context())

		// Extract method and service name
		methodName := getMethodName(info.FullMethod)
		serviceName := getServiceName(info.FullMethod)

		// Proceed with the original request
		err := handler(srv, stream)

		// Log details
		logger.Info("stream request",
			zap.String("protocol", "grpc"),
			zap.String("service", serviceName),
			zap.String("method", methodName),
			zap.String("method_type", "stream"),
			zap.String("ip", clientIP),
			zap.String("x-request-id", requestID),
			zap.Bool("client_streaming", info.IsClientStream),
			zap.Bool("server_streaming", info.IsServerStream),
			zap.Duration("duration", time.Since(startTime)),
			// Any other metadata you'd like to log
		)

		return err
	}
}

func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String()
		}
	}
	return ""
}

func getMethodName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullMethod
}

func getServiceName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) > 1 {
		// Further split to account for the package.ServiceName format
		subParts := strings.Split(parts[len(parts)-2], ".")
		return subParts[len(subParts)-1] // Return the last segment which is the service name
	}
	return ""
}
