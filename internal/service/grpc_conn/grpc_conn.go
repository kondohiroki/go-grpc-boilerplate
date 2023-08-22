package grpc_conn

import "google.golang.org/grpc"

// setup connection to gRPC server and return connection
// then other service can use this connection to call gRPC server with no need to setup connection again

var grpcConn *grpc.ClientConn
