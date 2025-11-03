package types

type Page struct {
	Data any  `json:"data"`
	More bool `json:"more"`
}

type Pagination struct {
	Page int
	Limit int
}
