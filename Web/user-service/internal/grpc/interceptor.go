package grpcserver

import (
	"context"
	"log"
	"strings"

	"Web/user-service/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwt *auth.JWTManager
}

func NewAuthInterceptor(jwt *auth.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwt: jwt}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/user.UserService/Register" ||
			info.FullMethod == "/user.UserService/Login" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		token := values[0]

		token = strings.TrimPrefix(token, "Bearer ")

		userID, err := a.jwt.Parse(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, "user_id", userID)
		log.Println("METHOD:", info.FullMethod)
		return handler(ctx, req)

	}

}
