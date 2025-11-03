package types

type RecommendationStore interface {
	InsertRecommendation(placeId int, userId int, reason string) error
}

type RecommendationPayload struct {
	Name string    `json:"name"`
	Address string `json:"address"`
}
