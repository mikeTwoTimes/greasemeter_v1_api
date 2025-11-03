package types

type PlaceStore interface {
	MakeEnvelope(box Bounds) ([]Marker, error)
	SuggestPlacesByLocation(term string, lat float64, lng float64) ([]SearchResult, error)
	GetPlaceListByLocation(box Bounds, page Pagination) (Page, error)
	GetMetaByPlace(placeId int) (PlaceMeta, error)
	GetInfoByPlace(placeId int) (PlaceInfo, error)
	GetImagesByPlace(placeId int) ([]string, error)
}

type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

type Marker struct {
	Id    int          `json:"id"`
	Point GeoJSONPoint `json:"point"`
}

type SearchResult struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type PlaceMeta struct {
	Id      int     `json:"id,omitempty"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Rating  float32 `json:"rating"`
}

type PlaceInfo struct {
	Rating float32  `json:"rating"`
	Images []string `json:"images"`
}

type Bounds struct {
	LatMin float64
	LatMax float64
	LngMin float64
	LngMax float64
}
