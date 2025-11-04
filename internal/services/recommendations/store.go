package recommendations

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateRecommendation(data types.RecommendationPayload, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        INSERT INTO recommendations (user_id, name, address)
        VALUES ($1, $2, $3)
    `
	
	_, err := s.db.Exec(
		ctx,
		query,
		userId,
		data.Name,
		data.Address,
	)

	return err
}
