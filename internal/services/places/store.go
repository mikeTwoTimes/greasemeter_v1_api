package places

import (
	"context"
	"encoding/json"
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

func (s *Store) GetMapMarkers(box types.Bounds) ([]types.Marker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            id,
            ST_AsGeoJSON(point)
        FROM places
        WHERE ST_WITHIN(point, ST_MakeEnvelope($1, $2, $3, $4, 4326))
    `

	rows, err := s.db.Query(
		ctx,
		query,
		box.LngMin,
		box.LatMin,
		box.LngMax,
		box.LatMax,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var markers []types.Marker

	for rows.Next() {
		var marker types.Marker
		geojson := ""
		err = rows.Scan(&marker.Id, &geojson)

		if err != nil {
			return nil, err
		} else if json.Unmarshal([]byte(geojson), &marker.Point) != nil {
			return nil, err
		}

		markers = append(markers, marker)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return markers, nil
}

func (s *Store) SearchForPlaces(term string, lat, lng float64) ([]types.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            id,
            name,
            address
        FROM places
        WHERE name ILIKE $1 || '%'
        ORDER BY ST_Distance(point, ST_MakePoint($2, $3)::geography) LIMIT 10
    `

	rows, err := s.db.Query(ctx, query, term, lat, lng)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var searchResults []types.SearchResult

	for rows.Next() {
		var searchResult types.SearchResult
		err = rows.Scan(
			&searchResult.Id,
			&searchResult.Name,
			&searchResult.Address,
		)

		if err != nil {
			return nil, err
		}

		searchResults = append(searchResults, searchResult)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return searchResults, nil
}

func (s *Store) GetPlacesList(box types.Bounds, page types.Pagination) (types.Page[types.PlaceMeta], error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            id,
            name,
            address,
            rating_sum,
            rating_count
        FROM places
        WHERE ST_WITHIN(point, ST_MakeEnvelope($1, $2, $3, $4, 4326))
        ORDER BY rating_count DESC
        LIMIT $5 + 1 OFFSET ($6 - 1) * $5
    `

	rows, err := s.db.Query(
		ctx,
		query,
		box.LngMin,
		box.LatMin,
		box.LngMax,
		box.LatMax,
		page.Limit,
		page.Offset,
	)

	if err != nil {
		return types.Page[types.PlaceMeta]{}, err
	}

	defer rows.Close()
	var list []types.PlaceMeta

	for rows.Next() {
		var meta types.PlaceMeta
		sum, count := 0, 0
		err = rows.Scan(
			&meta.Id,
			&meta.Name,
			&meta.Address,
			&sum,
			&count,
		)

		if err != nil {
			return types.Page[types.PlaceMeta]{}, err
		} else if count == 0 {
			meta.Rating = 0
		} else {
			meta.Rating = float32(sum) / float32(count)
		}

		list = append(list, meta)
	}

	if err = rows.Err(); err != nil {
		return types.Page[types.PlaceMeta]{}, err
	} else if len(list) <= page.Limit {
		return types.Page[types.PlaceMeta]{
			Data: list,
			More: false,
		}, nil
	}

	return types.Page[types.PlaceMeta]{
		Data: list[:page.Limit],
		More: true,
	}, nil
}

func (s *Store) GetMetaForPlace(placeId int) (types.PlaceMeta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            name,
            address,
            rating_sum,
            rating_count
        FROM places
        WHERE id = $1
    `

	var meta types.PlaceMeta
	sum, count := 0, 0
	err := s.db.QueryRow(ctx, query, placeId).Scan(
		&meta.Name,
		&meta.Address,
		&sum,
		&count,
	)

	if err != nil {
		return types.PlaceMeta{}, err
	} else if count == 0 {
		meta.Rating = 0
	} else {
		meta.Rating = float32(sum) / float32(count)
	}

	return meta, nil
}

func (s *Store) GetInfoForPlace(placeId int) (types.PlaceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            rating_sum,
            rating_count,
            image_url
        FROM places
        WHERE id = $1
    `

	var info types.PlaceInfo
	sum, count := 0, 0
	err := s.db.QueryRow(ctx, query, placeId).Scan(
		&sum,
		&count,
		&info.Images,
	)

	if err != nil {
		return types.PlaceInfo{}, err
	} else if count == 0 {
		info.Rating = 0
	} else {
		info.Rating = float32(sum) / float32(count)
	}

	return info, nil
}

func (s *Store) GetImagesForPlace(placeId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT image_url
        FROM places
        WHERE id = $1
    `

	var images []string
	err := s.db.QueryRow(ctx, query, placeId).Scan(&images)

	if err != nil {
		return nil, err
	}

	return images, nil
}
