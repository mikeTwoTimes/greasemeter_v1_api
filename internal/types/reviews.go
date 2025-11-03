package types

import "time"

type ReviewStore interface {
	InsertReview(data ReviewPayload, userId int, placeId int) (Timestamp, error)
	GetReviewKeysAndRating(reviewId int) (ReviewRef, error)
	GetReviewsByPlace(placeId int, page Pagination, query string) (Page, error)
	GetReviewsByUser(userId int, page Pagination, query string) (Page, error)
	UpdateReview(data ReviewPayload, reviewId int, placeId int, diff int) (Timestamp, error)
	DeleteReview(reviewId int, placeId int, rating int) error
	
	getReviews(foreignId int, page Pagination, query string) (Page, error)
}

type Review struct {
	Id     int       `json:"id"`
    Name   string    `json:"name"`
	Rating int       `json:"rating"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

type Timestamp struct {
	Id int       `json:"id,omitempty"`
	T  time.Time `json:"timestamp"`
}

type ReviewPayload struct {
	Rating  int    `json:"rating"`
	Text    string `json:"text"`
}

type ReviewRef struct {
	UserId  int
	PlaceId int
	Rating  int
}
