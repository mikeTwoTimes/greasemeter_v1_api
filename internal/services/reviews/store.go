package reviews

import (
	"context"
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

func (s *Store) CreateReview(data types.ReviewPayload, placeId, userId int) (types.Timestamp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	tx, err := s.db.Begin(ctx)
	defer cancel()

	if err != nil {
		return types.Timestamp{}, err
	}

	defer tx.Rollback(ctx)

	insertReviewQuery := `
        INSERT INTO reviews (place_id, user_id, rating, text)
        VALUES ($1, $2, $3, $4)
        RETURNING id, date
    `
	
	var timestamp types.Timestamp
	err = tx.QueryRow(
		ctx,
		insertReviewQuery,
		placeId,
		userId,
		data.Rating,
		data.Text,
	).Scan(&timestamp.Id, &timestamp.Time)

	if err != nil {
		return types.Timestamp{}, err
	}

	updateRatingQuery := `
        UPDATE places SET
            rating_count = rating_count + 1,
            rating_sum = rating_sum + $1
        WHERE id = $2
    `
	
	_, err = tx.Exec(ctx, updateRatingQuery, data.Rating, placeId)

	if err != nil {
		return types.Timestamp{}, err
	} else if err = tx.Commit(ctx); err != nil {
		return types.Timestamp{}, err
	}

	return timestamp, nil
}

func (s *Store) GetReviews(query string, foreignId int, page types.Pagination) (types.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	rows, err := s.db.Query(ctx, query, foreignId, page.Limit, page.Offset)

	if err != nil {
		return types.Page{}, err
	}

	defer rows.Close()
	var reviews []types.Review

	for rows.Next() {
		var review types.Review
		err = rows.Scan(
			&review.Id,
			&review.Name,
			&review.Rating,
			&review.Text,
			&review.Time,
		)

		if err != nil {
			return types.Page{}, err
		}

		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return types.Page{}, err
	} else if len(reviews) <= page.Limit {
		return types.Page{
			Data: reviews,
			More: false,
		}, nil
	}

	return types.Page{
		Data: reviews[:page.Limit],
		More: true,
	}, nil
}

func (s *Store) GetReviewKeysAndRating(reviewId int) (types.ReviewRef, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `
        SELECT
            user_id,
            place_id,
            rating
        FROM reviews
        WHERE id = $1
    `

	var refs types.ReviewRef
	err := s.db.QueryRow(ctx, query, reviewId).Scan(
		&refs.UserId,
		&refs.PlaceId,
		&refs.Rating,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return types.ReviewRef{}, nil
		}

		return types.ReviewRef{}, err
	}

	return refs, nil
}

func (s *Store) GetReviewsForPlace(placeId int, page types.Pagination) (types.Page, error) {
	query := `
        SELECT
            r.id,
            u.name AS username,
            r.rating,
            r.text,
            r.date
        FROM reviews r 
        JOIN users u ON r.user_id = u.id 
        WHERE r.place_id = $1
        ORDER BY r.date DESC, r.id DESC
        LIMIT $2 + 1 OFFSET ($3 - 1) * $2
    `
	
	 return s.GetReviews(query, placeId, page)
}

func (s *Store) GetReviewsForUser(userId int, page types.Pagination) (types.Page, error) {
	query := `
        SELECT
            r.id,
            p.name AS place_name,
            r.rating,
            r.text,
            r.date
        FROM reviews r 
        JOIN places p ON r.place_id = p.id 
        WHERE r.user_id = $1
        ORDER BY r.date DESC, r.id DESC
        LIMIT $2 + 1 OFFSET ($3 - 1) * $2
    `
	
	return s.GetReviews(query, userId, page)
}

func (s *Store) UpdateReview(data types.ReviewPayload, reviewId, placeId, diff int) (types.Timestamp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	tx, err := s.db.Begin(ctx)
	defer cancel()

	if err != nil {
		return types.Timestamp{}, err
	}

	defer tx.Rollback(ctx)

	updateReviewQuery := `
        UPDATE reviews SET
            rating = $1,
            text = $2
        WHERE id = $3
        RETURNING date
    `
	
	var timestamp types.Timestamp
	err = tx.QueryRow(
		ctx,
		updateReviewQuery,
		data.Rating,
		data.Text,
		reviewId,
	).Scan(&timestamp.Time)

	if err != nil {
		return types.Timestamp{}, err
	}

	updateRatingQuery := `
        UPDATE places SET
            rating_sum = rating_sum + $1
        WHERE id = $2
    `

	_, err = tx.Exec(ctx, updateRatingQuery, diff, placeId)

	if err != nil {
		return types.Timestamp{}, err
	} else if err = tx.Commit(ctx); err != nil {
		return types.Timestamp{}, err
	}

	return timestamp, nil
}

func (s *Store) DeleteReview(reviewId, placeId, rating int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	tx, err := s.db.Begin(ctx)
	defer cancel()

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	deleteReviewQuery := `
        DELETE FROM reviews
        WHERE id = $1
    `

	if _, err = tx.Exec(ctx, deleteReviewQuery, reviewId); err != nil {
		return err
	}

	updateRatingQuery := `
        UPDATE places SET
            rating_count = rating_count - 1,
            rating_sum = rating_sum - $1
        WHERE id = $2
    `

	_, err = tx.Exec(ctx, updateRatingQuery, rating, placeId)

	if err != nil {
		return err
	} else if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
