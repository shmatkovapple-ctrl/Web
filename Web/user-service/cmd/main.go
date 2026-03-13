package main

import (
	userpb "Web/user-service/protos/gen/go"
	"log"
	"net"

	"Web/user-service/internal/auth"
	"Web/user-service/internal/config"
	"Web/user-service/internal/database"
	"Web/user-service/internal/grpc"
	"Web/user-service/internal/repository"
	"Web/user-service/internal/service"

	"google.golang.org/grpc"
)

func main() {

	cfg := config.LoadConfig()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(db)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userService := service.NewUserService(repo, jwtManager)

	authInterceptor := grpcserver.NewAuthInterceptor(jwtManager)

	log.Println("PORT =", cfg.GRPCPort)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	userHandler := grpcserver.NewServer(userService)

	userpb.RegisterUserServiceServer(
		grpcServer,
		userHandler,
	)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC started on port", cfg.GRPCPort)
	grpcServer.Serve(lis)
}
