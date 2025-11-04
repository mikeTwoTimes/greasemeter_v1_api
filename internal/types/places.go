package types

type PlaceStore interface {
	GetMapMarkers(box Bounds) ([]Marker, error)
	SearchForPlaces(term string, lat, lng float64) ([]SearchResult, error)
	GetPlacesList(box Bounds, page Pagination) (Page, error)
	GetMetaForPlace(placeId int) (PlaceMeta, error)
	GetInfoForPlace(placeId int) (PlaceInfo, error)
	GetImagesForPlace(placeId int) ([]string, error)
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
