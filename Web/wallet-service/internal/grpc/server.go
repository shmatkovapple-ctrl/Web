package grpc

import (
	"Web/wallet-service/internal/repository"
	walletpb "Web/wallet-service/protos/gen/go"
	"context"
	"log"
)

type Server struct {
	walletpb.UnimplementedWalletServiceServer
	service repository.WalletService
}

func NewServer(s repository.WalletService) *Server {
	return &Server{
		service: s,
	}
}

func (s *Server) CreateWallet(ctx context.Context, req *walletpb.CreateWalletRequest) (*walletpb.CreateWalletResponse, error) {
	_, err := s.service.CreateWallet(
		ctx,
		int(req.UserId),
		req.Currency,
	)
	if err != nil {
		log.Println(err)
	}
	return &walletpb.CreateWalletResponse{
		UserID: int(req.UserId),
	}, nil
}

func (s *Server) GetBalance(
	ctx context.Context,
	req *walletpb.GetBalanceRequest,
) (*walletpb.GetBalanceResponse, error) {

	balance, err := s.service.GetBalance(
		ctx,
		int(req.UserId),
	)

	if err != nil {
		return nil, err
	}

	return &walletpb.GetBalanceResponse{
		Balance: balance,
	}, nil
}

func (s *Server) Deposit(
	ctx context.Context,
	req *walletpb.DepositRequest,
) (*walletpb.DepositResponse, error) {

	balance, err := s.service.Deposit(
		ctx,
		int(req.UserId),
		req.Amount,
	)

	if err != nil {
		return nil, err
	}

	return &walletpb.DepositResponse{
		Balance: balance,
	}, nil
}
