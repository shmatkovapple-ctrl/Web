package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: pool}
}

func (r *PostgresRepository) CreateWallet(ctx context.Context, userID int, currency string) (int, error) {
	var id int

	err := r.db.QueryRow(ctx,
		`INSERT INTO wallet (user_id, currency, balance) VALUES ($1, $2, 0) RETURNING wallet_id`,
		userID, currency).Scan(&id)
	if err != nil {
		log.Println("Ошибка создания кошелька")
	}

	log.Println("Кошелёк создан:", id)

	return id, nil
}

func (r *PostgresRepository) GetBalance(ctx context.Context, id int) (int64, error) {

	var balance int

	err := r.db.QueryRow(ctx,
		`SELECT balance FROM wallet WHERE user_id = $1`, id).Scan(&balance)
	if err != nil {
		log.Println("Ошибка чтения кошелька")
	}
	return int64(balance), err
}

func (r *PostgresRepository) Deposit(ctx context.Context, userID int, amount int64) (int64, error) {
	var balance int64

	err := r.db.QueryRow(ctx,
		`UPDATE wallet SET balance = balance + $1 where user_id = $2 RETURNING balance`, amount, userID).Scan(&balance)
	if err != nil {
		log.Println("Ошибка депозита")
	}
	return balance, err
}
