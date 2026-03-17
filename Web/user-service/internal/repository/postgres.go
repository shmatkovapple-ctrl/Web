package repository

import (
	"Web/user-service/internal/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, u *models.User) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (login, first_name, last_name, birth_date, password_hash) 
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		u.Login, u.FirstName, u.LastName, u.BirthDate, u.PasswordHash).Scan(&id)
	if err != nil {
		log.Println("DB ERROR:", err)
	}

	log.Println("USER CREATED ID:", id)

	return id, err
}

func (r *PostgresRepo) GetByFirstName(ctx context.Context, login string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, login, first_name, last_name, birth_date, password_hash
			 FROM users WHERE login = $1`, login)

	u := &models.User{}
	err := row.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.BirthDate, &u.PasswordHash)
	if err != nil {
		log.Println("DB ERROR:", err)
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepo) GetById(ctx context.Context, id int) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, login, first_name, last_name, birth_date
			 FROM users WHERE id = $1`, id)

	u := &models.User{}
	err := row.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.BirthDate)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
