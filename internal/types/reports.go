package types

type ReportStore interface {
	InsertReport(placeId int, userId int, reason string) error
}

type ReportPayload struct {
	Reason string `json:"reason"`
}
