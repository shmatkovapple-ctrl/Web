package service

import (
	walletpb "Web/user-service/protos/gen/go/wallet"
	"context"

	"google.golang.org/grpc"
)

type WalletService struct {
	client walletpb.WalletServiceClient
}

func NewWalletService(conn *grpc.ClientConn) *WalletService {
	return &WalletService{
		client: walletpb.NewWalletServiceClient(conn),
	}
}

func (c *WalletService) GetBalance(ctx context.Context, userID int32) (int64, error) {
	resp, err := c.client.GetBalance(ctx, &walletpb.GetBalanceRequest{
		UserId: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.Balance, nil
}

func (c *WalletService) CreateWallet(ctx context.Context, userID int32) error {
	_, err := c.client.CreateWallet(ctx, &walletpb.CreateWalletRequest{
		UserId:   userID,
		Currency: "USD",
	})
	return err
}
