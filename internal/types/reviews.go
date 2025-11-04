package types

import "time"

type ReviewStore interface {
	CreateReview(data ReviewPayload, userId, placeId int) (Timestamp, error)
	GetReviewKeysAndRating(reviewId int) (ReviewRef, error)
	GetReviewsForPlace(placeId int, page Pagination) (Page, error)
	GetReviewsForUser(userId int, page Pagination) (Page, error)
	GetReviews(query string, foreignId int, page Pagination) (Page, error)
	UpdateReview(data ReviewPayload, reviewId, placeId, diff int) (Timestamp, error)
	DeleteReview(reviewId, placeId, rating int) error
}

type Review struct {
	Id     int       `json:"id"`
    Name   string    `json:"name"`
	Rating int       `json:"rating"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}

type Timestamp struct {
	Id    int       `json:"id,omitempty"`
	Time  time.Time `json:"time"`
}

type ReviewPayload struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

type ReviewRef struct {
	UserId  int
	PlaceId int
	Rating  int
}
