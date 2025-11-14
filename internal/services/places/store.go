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

func (s *Store) GetPlacesList(box types.Bounds, page types.Pagination) (types.Page[types.Listing], error) {
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
		return types.Page[types.Listing]{}, err
	}

	defer rows.Close()
	var list []types.Listing

	for rows.Next() {
		var place types.Listing
		sum, count := 0, 0
		err = rows.Scan(
			&place.Id,
			&place.Name,
			&place.Address,
			&sum,
			&count,
		)

		if err != nil {
			return types.Page[types.Listing]{}, err
		} else if count == 0 {
			place.Rating = 0
		} else {
			place.Rating = float32(sum) / float32(count)
		}

		list = append(list, place)
	}

	if err = rows.Err(); err != nil {
		return types.Page[types.Listing]{}, err
	} else if len(list) <= page.Limit {
		return types.Page[types.Listing]{
			Data: list,
			More: false,
		}, nil
	}

	return types.Page[types.Listing]{
		Data: list[:page.Limit],
		More: true,
	}, nil
}

func (s *Store) GetMarkerDetails(placeId int) (types.MarkerDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            name,
            address,
            rating_sum,
            rating_count,
            image_url
        FROM places
        WHERE id = $1
    `

	var details types.MarkerDetails
	sum, count := 0, 0
	err := s.db.QueryRow(ctx, query, placeId).Scan(
		&details.Name,
		&details.Address,
		&sum,
		&count,
		&details.Images,
	)

	if err != nil {
		return types.MarkerDetails{}, err
	} else if count == 0 {
		details.Rating = 0
	} else {
		details.Rating = float32(sum) / float32(count)
	}

	return details, nil
}

func (s *Store) GetListingDetails(placeId int) (types.ListingDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            ST_AsGeoJSON(point),
            image_url
        FROM places
        WHERE id = $1
    `

	var details types.ListingDetails
	geojson := ""
	err := s.db.QueryRow(ctx, query, placeId).Scan(
		&geojson,
		&details.Images,
	)

	if err != nil {
		return types.ListingDetails{}, err
	} else if json.Unmarshal([]byte(geojson), &details.Point) != nil {
		return types.ListingDetails{}, err
	}

	return details, nil
}

func (s *Store) GetPlaceMeta(placeId int) (types.PlaceMeta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT
            ST_AsGeoJSON(point),
            rating_sum,
            rating_count,
            image_url
        FROM places
        WHERE id = $1
    `

	var meta types.PlaceMeta
	geojson := ""
	sum, count := 0, 0
	err := s.db.QueryRow(ctx, query, placeId).Scan(
		&geojson,
		&sum,
		&count,
		&meta.Images,
	)

	if err != nil {
		return types.PlaceMeta{}, err
	} else if json.Unmarshal([]byte(geojson), &meta.Point) != nil {
		return types.PlaceMeta{}, err
	} else if count == 0 {
		meta.Rating = 0
	} else {
		meta.Rating = float32(sum) / float32(count)
	}

	return meta, nil
}
