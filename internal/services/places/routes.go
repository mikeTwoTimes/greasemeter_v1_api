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

func (h *Handler) getMetaForPlace(c *gin.Context) {
	getPlaceData(c, h.store.GetMetaForPlace)
}

func (h *Handler) getInfoForPlace(c *gin.Context) {
	getPlaceData(c, h.store.GetInfoForPlace)
}

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
