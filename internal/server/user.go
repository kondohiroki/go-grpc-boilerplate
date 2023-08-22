package server

import (
	"context"
	"fmt"

	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
)

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	dto, err := s.app.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Replace with actual logic to retrieve the user from the database.
	return &pb.GetUserResponse{
		Status: &pb.Status{
			Code:    0,
			Message: "success",
		},
		Data: &pb.User{
			FirstName: dto[0].Name,
			LastName:  dto[0].Email,
			Point:     100,
		},
	}, nil
}

func (s *Server) GetUserList(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserListResponse, error) {
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

func (s *Server) GetUserPagination(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserPaginationResponse, error) {
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

func (s *Server) GetGrpcError(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	logger.Log.Info("GetGrpcError invoked")

	return &pb.GetUserResponse{
		Status: &pb.Status{
			Code:    100,
			Message: "error",
		},
		Data: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Point:     100,
		},
	}, nil
}
