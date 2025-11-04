package types

type ReportStore interface {
	CreateReport(placeId, userId int, reason string) error
}

type ReportPayload struct {
	Reason string `json:"reason"`
}
