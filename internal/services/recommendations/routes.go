package recommendations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
)

type Handler struct {
	store types.RecommendationStore
}

func NewHandler(store types.RecommendationStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.POST("/recommendations", h.createRecommendation)
}

// @Summary	    Creates a place recommendation
// @Description	Creates a place recommendation given a valid access token
// @Tags        recommendations
// @Accept      json
// @Produce     json
// @Param       reccommendation body types.RecommendationPayload true "Recommendation"
// @Success	    204
// @Router      /v1/recommendations [post]
// @Security    BearerAuth
func (h *Handler) createRecommendation(c *gin.Context) {
	req, err := utility.ParseRecommendation(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(int)

	if err = h.store.CreateRecommendation(req, userId); err != nil {
		c.JSON(utility.MapError(err))
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
