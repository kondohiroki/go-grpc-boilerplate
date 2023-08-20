package server

import (
	"context"

	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	pb "github.com/kondohiroki/go-grpc-boilerplate/proto"
)

func (s *Server) GetUser(context.Context, *pb.GetUserReq) (*pb.GetUserReply, error) {

	logger.Log.Info("GetUsers invoked")

	// Replace with actual logic to retrieve the user from the database.
	return &pb.GetUserReply{
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

/*
	GetUserList(context.Context, *pb.GetUserReq) (*pb.GetUserListReply, error)
	GetUserPagination(context.Context, *pb.GetUserReq) (*pb.GetUserPaginationReply, error)
	GetGrpcError(context.Context, *pb.GetUserReq) (*pb.GetUserReply, error)
*/

func (s *Server) GetUserList(context.Context, *pb.GetUserReq) (*pb.GetUserListReply, error) {
	logger.Log.Info("GetUserList invoked")

	return &pb.GetUserListReply{
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

func (s *Server) GetUserPagination(context.Context, *pb.GetUserReq) (*pb.GetUserPaginationReply, error) {
	logger.Log.Info("GetUserPagination invoked")

	return &pb.GetUserPaginationReply{
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

// func (s *Server) GetGrpcError(context.Context, *pb.GetUserReq) (*pb.GetUserReply, error) {
// 	logger.Log.Info("GetGrpcError invoked")

// 	return &pb.GetUserReply{
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
