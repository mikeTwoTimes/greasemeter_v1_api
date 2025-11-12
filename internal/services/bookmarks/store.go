package bookmarks

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

func (s *Store) CreateBookmark(userId, placeId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO bookmarks (user_id, place_id)
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(ctx, query, userId, placeId)

	return err
}

func (s *Store) GetUserFromBookmark(bookmarkId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT user_id
        FROM bookmarks
        WHERE id = $1
    `

	var userId int
	err := s.db.QueryRow(ctx, query, bookmarkId).Scan(&userId)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return userId, nil
}

func (s *Store) GetBookmarksForUser(userId int) ([]types.Bookmark, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            b.id AS bookmark_id,
            b.place_id,
            p.name AS place_name,
            p.address
        FROM bookmarks b
        JOIN places p ON b.place_id = p.id
        WHERE b.user_id = $1
        ORDER BY b.id
    `

	rows, err := s.db.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var bookmarks []types.Bookmark

	for rows.Next() {
		var bookmark types.Bookmark
		err = rows.Scan(
			&bookmark.Id,
			&bookmark.PlaceId,
			&bookmark.Name,
			&bookmark.Address,
		)

		if err != nil {
			return nil, err
		}

		bookmarks = append(bookmarks, bookmark)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (s *Store) IsPlaceBookmarked(userId, placeId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM bookmarks
            WHERE user_id = $1 AND place_id = $2
        )
    `

	var exists bool
	err := s.db.QueryRow(ctx, query, userId, placeId).Scan(&exists)

	return exists, err
}

func (s *Store) DeleteBookmark(bookmarkId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        DELETE FROM bookmarks
        WHERE id = $1
    `

	_, err := s.db.Exec(ctx, query, bookmarkId)

	return err
}
