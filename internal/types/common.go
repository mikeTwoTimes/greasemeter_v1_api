package types

type Page[T any] struct {
	Data []T  `json:"data"`
	More bool `json:"more"`
}

type Pagination struct {
	Offset int
	Limit  int
}
