package server

import (
	"context"

	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
)

func (s *Server) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Replace with actual logic to retrieve the user from the database.
	return &pb.GetUserResponse{
		Status: &pb.Status{
			Code:    0,
			Message: "success",
		},
		Data: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Point:     100,
		},
	}, nil
}

func (s *Server) GetUserList(context.Context, *pb.GetUserRequest) (*pb.GetUserListResponse, error) {
	logger.Log.Info("GetUserList invoked")

	return &pb.GetUserListResponse{
		Data: []*pb.User{
			{
				FirstName: "John",
				LastName:  "Doe",
				Point:     100,
			},
			{
				FirstName: "Jane",
				LastName:  "Doe",
				Point:     200,
			},
		},
	}, nil
}

func (s *Server) GetUserPagination(context.Context, *pb.GetUserRequest) (*pb.GetUserPaginationResponse, error) {
	logger.Log.Info("GetUserPagination invoked")

	return &pb.GetUserPaginationResponse{
		Data: []*pb.User{
			{
				FirstName: "John",
				LastName:  "Doe",
				Point:     100,
			},
			{
				FirstName: "Jane",
				LastName:  "Doe",
				Point:     200,
			},
		},
	}, nil
}

// func (s *Server) GetGrpcError(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
// 	logger.Log.Info("GetGrpcError invoked")

// 	return &pb.GetUserResponse{
// 		Status: &pb.Status{
// 			Code:    100,
// 			Message: "error",
// 		},
// 		Data: &pb.User{
// 			FirstName: "John",
// 			LastName:  "Doe",
// 			Point:     100,
// 		},
// 	}, nil
// }
