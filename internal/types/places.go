package types

type PlaceStore interface {
	GetMapMarkers(box Bounds) ([]Marker, error)
	SearchForPlaces(term string, lat, lng float64) ([]SearchResult, error)
	GetPlacesList(box Bounds, page Pagination) (Page[Listing], error)
	GetMarkerDetails(placeId int) (MarkerDetails, error)
	GetListingDetails(placeId int) (ListingDetails, error)
	GetPlaceMeta(placeId int) (PlaceMeta, error)
}

type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

type Marker struct {
	Id    int          `json:"id"`
	Point GeoJSONPoint `json:"point"`
}

type MarkerDetails struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Rating  float32  `json:"rating"`
	Images  []string `json:"images"`
}

type Listing struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Rating  float32 `json:"rating"`
}

type ListingDetails struct {
	Point  GeoJSONPoint `json:"point"`
	Images []string     `json:"images"`
}

type SearchResult struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type PlaceMeta struct {
	Point  GeoJSONPoint `json:"point"`
	Rating float32      `json:"rating"`
	Images []string     `json:"images"`
}

type Bounds struct {
	LatMin float64
	LatMax float64
	LngMin float64
	LngMax float64
}

type ListingPage struct {
	Data []Listing `json:"data"`
	More bool      `json:"more"`
}
