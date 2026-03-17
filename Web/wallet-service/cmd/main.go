package main

import (
	"Web/wallet-service/internal/config"
	"Web/wallet-service/internal/database"
	grpcserver "Web/wallet-service/internal/grpc"
	"Web/wallet-service/internal/repository"
	walletpb "Web/wallet-service/protos/gen/go"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db, _ := database.NewPostgres(cfg)

	repo := repository.NewPostgresRepository(db)

	walletService := repository.NewWalletService(repo)

	grpcServer := grpc.NewServer()

	handler := grpcserver.NewServer(*walletService)

	walletpb.RegisterWalletServiceServer(
		grpcServer,
		handler,
	)

	lis, _ := net.Listen("tcp", ":"+cfg.GRPCPort)

	log.Println("gRPC started on port", cfg.GRPCPort)

	grpcServer.Serve(lis)
}
