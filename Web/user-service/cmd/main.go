package main

import (
	"Web/user-service/internal/service"
	"Web/user-service/internal/usecase"
	userpb "Web/user-service/protos/gen/go"
	"log"
	"net"

	"Web/user-service/internal/auth"
	"Web/user-service/internal/config"
	"Web/user-service/internal/database"
	"Web/user-service/internal/grpc"
	"Web/user-service/internal/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg := config.LoadConfig()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(db)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	conn, err := grpc.Dial(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	walletService := service.NewWalletService(conn)

	userService := usecase.NewUserService(repo, jwtManager, walletService)

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
