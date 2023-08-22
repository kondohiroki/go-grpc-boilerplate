package middleware

import (
	"context"
	"fmt"

	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// requestIDKey is a custom type for context keys used for the request ID
type requestIDKey struct{}

func getRequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey{}).(string); ok {
		return reqID
	}
	return "unknown"
}

func addCommonHeaders(requestID string) metadata.MD {
	return metadata.Pairs(
		"x-request-id", requestID,
		"x-received-time", carbon.Now().ToString(),
		"x-app-name", config.GetConfig().App.NameSlug,
		"x-app-version", "unknown",
		"x-app-unix-time", fmt.Sprintf("%d", carbon.Now().Timestamp()),
	)
}

// UnaryRequestIDInterceptor is a gRPC server-side interceptor that adds the request ID to the context.
func UnaryRequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	requestIDs, ok := md["x-request-id"]
	var requestID string
	if !ok || len(requestIDs) == 0 {
		requestID = uuid.New().String() // Generate a new ID if not provided
	} else {
		requestID = requestIDs[0]
	}

	ctx = context.WithValue(ctx, requestIDKey{}, requestID)

	resp, err := handler(ctx, req)

	headers := addCommonHeaders(requestID)
	grpc.SetHeader(ctx, headers)

	return resp, err
}

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying the Context.
type WrappedServerStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a new ServerStream that has the ability to modify the context.
func WrapServerStream(ss grpc.ServerStream) *WrappedServerStream {
	return &WrappedServerStream{ServerStream: ss, WrappedContext: ss.Context()}
}

// StreamRequestIDInterceptor is a gRPC server-side interceptor that adds the request ID to the context.
func StreamRequestIDInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	md, _ := metadata.FromIncomingContext(ctx)
	requestIDs, ok := md["x-request-id"]
	var requestID string
	if !ok || len(requestIDs) == 0 {
		requestID = uuid.New().String() // Generate a new ID if not provided
	} else {
		requestID = requestIDs[0]
	}

	headers := addCommonHeaders(requestID)
	ss.SendHeader(headers)

	ctx = context.WithValue(ctx, requestIDKey{}, requestID)
	wrappedStream := WrapServerStream(ss)
	wrappedStream.WrappedContext = ctx

	return handler(srv, wrappedStream)
}

// StreamClientRequestIDInterceptor is a gRPC client-side interceptor that adds the request ID to the outgoing context.
func StreamClientRequestIDInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	requestID := uuid.New().String() // Generate a new ID
	md := metadata.Pairs("x-request-id", requestID)
	newCtx := metadata.NewOutgoingContext(ctx, md)

	// Call the original streamer with the new context
	return streamer(newCtx, desc, cc, method, opts...)
}
