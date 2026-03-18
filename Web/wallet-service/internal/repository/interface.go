package repository

import (
	"context"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, id int, currency string) (int, error)
	GetBalance(ctx context.Context, id int) (int64, error)
	Deposit(ctx context.Context, userID int, amount int64) (int64, error)
}

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) CreateWallet(ctx context.Context, userID int, currency string) (int, error) {
	return s.repo.CreateWallet(ctx, userID, currency)
}

func (s *WalletService) GetBalance(ctx context.Context, id int) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}

func (s *WalletService) Deposit(ctx context.Context, userID int, amount int64) (int64, error) {
	return s.repo.Deposit(ctx, userID, amount)
}
