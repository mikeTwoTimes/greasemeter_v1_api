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
