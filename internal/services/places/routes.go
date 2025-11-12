package places

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
)

type Handler struct {
	store types.PlaceStore
}

func NewHandler(store types.PlaceStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup) {
	v1.GET("/places/map", h.getMapMarkers)
	v1.GET("/places/search", h.searchForPlaces)
	v1.GET("/places/list", h.getPlacesList)
	v1.GET("/places/:id/meta", h.getMetaForPlace)
	v1.GET("/places/:id/info", h.getInfoForPlace)
	v1.GET("/places/:id/images", h.getImagesForPlace)
}

// @Summary	    Gets all map markers
// @Description	Gets all map markers within a specified region
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       lat query float64 true "Latitude"
// @Param       lng query float64 true "Longitude"
// @Param       latDelta query float64 true "Latitude delta"
// @Param       lngDelta query float64 true "Longitude delta"
// @Success	    200 {object} []types.Marker
// @Router      /v1/places/map [get]
func (h *Handler) getMapMarkers(c *gin.Context) {
	box, err := utility.ParseBoundingBox(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.store.GetMapMarkers(box)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Gets search results for places
// @Description	Gets search results given a search term, and point
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       term query string true "Search term"
// @Param       lat query float64 true "Latitude"
// @Param       lng query float64 true "Longitude"
// @Success	    200 {object} []types.SearchResult
// @Router      /v1/places/search [get]
func (h *Handler) searchForPlaces(c *gin.Context) {
	lat, lng, err := utility.ParseCoordinates(c)
	term := c.Query("term")

	if len([]rune(term)) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search term must be 1 character or more",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.store.SearchForPlaces(term, lat, lng)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Gets list of places
// @Description	Gets list of places within a specified region
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       lat query float64 true "Latitude"
// @Param       lng query float64 true "Longitude"
// @Param       latDelta query float64 true "Latitude delta"
// @Param       lngDelta query float64 true "Longitude delta"
// @Param       page query int true "Page number"
// @Param       limit query int true "Page length"
// @Success	    200 {object} types.MetaPage
// @Router      /v1/places/list [get]
func (h *Handler) getPlacesList(c *gin.Context) {
	box, err := utility.ParseBoundingBox(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := utility.ParsePagination(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.store.GetPlacesList(box, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Gets metadata for a place
// @Description	Gets metadata for a place given a valid place ID
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    200 {object} types.PlaceMeta
// @Router      /v1/places/{id}/meta [get]
func (h *Handler) getMetaForPlace(c *gin.Context) {
	getPlaceData(c, h.store.GetMetaForPlace)
}

// @Summary	    Gets info for a place
// @Description	Gets info for a place given a valid place ID
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    200 {object} types.PlaceInfo
// @Router      /v1/places/{id}/info [get]
func (h *Handler) getInfoForPlace(c *gin.Context) {
	getPlaceData(c, h.store.GetInfoForPlace)
}

// @Summary	    Gets images for a place
// @Description	Gets images for a place given a valid place ID
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    200 {array} string
// @Router      /v1/places/{id}/images [get]
func (h *Handler) getImagesForPlace(c *gin.Context) {
	getPlaceData(c, h.store.GetImagesForPlace)
}

func getPlaceData[T any](c *gin.Context, fetch func(int) (T, error)) {
	placeId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
		return
	}

	data, err := fetch(placeId)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, data)
	}
}
