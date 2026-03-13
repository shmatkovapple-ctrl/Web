package grpcserver

import (
	"Web/user-service/internal/service"
	"Web/user-service/protos/gen/go"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewServer(s *service.UserService) *Server {
	return &Server{service: s}
}

func (s *Server) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.AuthResponse, error) {

	token, err := s.service.Register(ctx,
		req.Login,
		req.FirstName,
		req.LastName,
		req.BirthDate,
		req.Password,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.AuthResponse{Token: token}, nil
}

func (s *Server) GetProfile(
	ctx context.Context,
	req *userpb.Empty,
) (*userpb.UserResponse, error) {

	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	user, err := s.service.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userpb.UserResponse{
		Login:     user.Login,
		Id:        int32(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}, nil
}

func (s *Server) Login(
	ctx context.Context,
	req *userpb.LoginRequest,
) (*userpb.AuthResponse, error) {

	token, err := s.service.Login(
		ctx,
		req.Login,
		req.Password,
	)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.AuthResponse{
		Token: token,
	}, nil
}
