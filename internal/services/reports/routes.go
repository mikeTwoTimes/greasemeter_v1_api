package reports

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
)

type Handler struct {
	store types.ReportStore
}

func NewHandler(store types.ReportStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.POST("/report/places/:id", h.createReport)
}

func (h *Handler) createReport(c *gin.Context) {
	placeId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
        return
    }

	reason, err := utility.ParseReport(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	userId := utility.GetUserFromContext(c)

	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	} else if err = h.store.CreateReport(placeId, userId, reason); err != nil {
		c.JSON(utility.MapError(err))
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
