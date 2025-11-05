package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(data types.RegisterPayload) error {
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

func (s *Store) CreateResetToken(userId int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return "", errors.New("Failed to generate token")
	}

	tokenString := hex.EncodeToString(b)

	query := `
        INSERT INTO reset_tokens (token, user_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id)
        DO UPDATE SET
            token = EXCLUDED.token,
            expires_at = NOW() + INTERVAL '15 minutes'
    `

	_, err = s.db.Exec(
		ctx,
		query,
		tokenString,
		userId,
	)

	if err != nil {
		return "", err
	}

	return tokenString, nil
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
			return types.Credentials{}, nil
		}

		return types.Credentials{}, err
	}

	return cred, nil
}

func (s *Store) GetUserFromEmail(email string) (int, error) {
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

func (s *Store) GetDataFromResetToken(token string) (types.ResetTokenData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        SELECT
            user_id,
            expires_at
        FROM reset_tokens
        WHERE token = $1
    `

	var data types.ResetTokenData
	err := s.db.QueryRow(ctx, query, token).Scan(
		&data.UserId,
		&data.Expiration,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return types.ResetTokenData{}, nil
		}

		return types.ResetTokenData{}, err
	}

	return data, nil
}

func (s *Store) UpdateUserPassword(userId int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	tx, err := s.db.Begin(ctx)
	defer cancel()

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	updatePasswordQuery := `
        UPDATE users SET
            password = $1
        WHERE id = $2
    `

	_, err = tx.Exec(ctx, updatePasswordQuery, password, userId)

	if err != nil {
		return err
	}

	deleteTokenQuery := `
        DELETE FROM reset_tokens
        WHERE user_id = $1
    `

	_, err = tx.Exec(ctx, deleteTokenQuery, userId)

	if err != nil {
		return err
	} else if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
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
