package reports

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateReport(placeId, userId int, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        INSERT INTO reports (user_id, place_id, reason)
        VALUES ($1, $2, $3)
    `

	_, err := s.db.Exec(
		ctx,
		query,
		placeId,
		userId,
		reason,
	)

	return err
}
