package users

import (
	"Greasemeter-rest-api/internal/types"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) InsertUser(data types.RegisterPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        INSERT INTO USERS (email, name, password)
        VALUES ($1, $2, $3)
    `

	_, err := s.db.Exec(
		ctx,
		query,
		data.Email,
		data.Name,
		data.Password,
	)

	return err
}

func (s *Store) GetUserCredentials(name string) (types.Credentials, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        SELECT
            id,
            password
        FROM users
        WHERE name = $1
    `

	var cred types.Credentials
	err := s.db.QueryRow(ctx, query, name).Scan(
		&cred.Id,
		&cred.Password,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return cred, nil
}

func (s *Store) GetUserByEmail(email string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        SELECT id
        FROM users
        WHERE email = $1
    `

	var id int
	err := s.db.QueryRow(ctx, query, email).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return id, nil
}

func (s *Store) UserExists(userId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM users
            WHERE id = $1
        )
    `

	var exists bool
	err := s.db.QueryRow(ctx, query, userId).Scan(&exists)

	return exists, err
}

func (s *Store) UpdateUserPassword(userId int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

	query := `
        UPDATE users SET
            password = $1
        WHERE id = $2
    `

	_, err := s.db.Exec(ctx, query, password, userId)

	return err
}

func (s *Store) DeleteUser(userId int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

	query := `
        DELETE FROM users
        WHERE id = $1
    `

	_, err := s.db.Exec(ctx, query, userId)

	return err
}
