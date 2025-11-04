package types

type RecommendationStore interface {
	CreateRecommendation(data RecommendationPayload, userId int) error
}

type RecommendationPayload struct {
	Name string    `json:"name"`
	Address string `json:"address"`
}
