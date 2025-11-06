package reviews

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
)

type Handler struct {
	store types.ReviewStore
}

func NewHandler(store types.ReviewStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(v1, auth *gin.RouterGroup) {
	v1.GET("/reviews/places/:id", h.getReviewsForPlace)

	auth.POST("/reviews/places/:id", h.createReview)
	auth.GET("/reviews", h.getReviewsForUser)
	auth.PATCH("/reviews/:id", h.updateReview)
	auth.DELETE("/reviews/:id", h.deleteReview)
}

// @Summary	    Creates a review for a place
// @Description	Creates a valid user's review for a place
// @Tags        reviews
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Param       review body types.ReviewPayload true "Review"
// @Success	    201 {object} types.Timestamp
// @Router      /v1/reviews/places/{id} [post]
// @Security    BearerAuth
func (h *Handler) createReview(c *gin.Context) {
	var req types.ReviewPayload
	placeId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
		return
	} else if req, err = utility.ParseReview(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(int)
	resp, err := h.store.CreateReview(req, placeId, userId)

	if err != nil {
		c.JSON(utility.MapError(err))
	} else {
		c.JSON(http.StatusCreated, resp)
	}
}

// @Summary	    Gets reviews for a place
// @Description	Gets reviews for a place given a valid place ID
// @Tags        reviews
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    200 {object} types.ReviewPage
// @Router      /v1/reviews/places/{id} [get]
func (h *Handler) getReviewsForPlace(c *gin.Context) {
	placeId, err := strconv.Atoi(c.Param("id"))
	var page types.Pagination

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
		return
	} else if page, err = utility.ParsePagination(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	resp, err := h.store.GetReviewsForPlace(placeId, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve reviews",
		})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Gets user's reviews
// @Description	Gets reviews for a user given a valid access token
// @Tags        reviews
// @Accept      json
// @Produce     json
// @Success	    200 {object} types.ReviewPage
// @Router      /v1/reviews/ [get]
// @Security    BearerAuth
func (h *Handler) getReviewsForUser(c *gin.Context) {
	page, err := utility.ParsePagination(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	userId := c.MustGet("userId").(int)
    resp, err := h.store.GetReviewsForUser(userId, page)
	
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Updates a users review
// @Description	Updates a users review given a valid access token
// @Tags        reviews
// @Accept      json
// @Produce     json
// @Param       id path int true "Review ID"
// @Success	    200 {object} types.Timestamp
// @Router      /v1/reviews/{id} [patch]
// @Security    BearerAuth
func (h *Handler) updateReview(c *gin.Context) {
	var req types.ReviewPayload
	reviewId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	} else if req, err = utility.ParseReview(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(int)
	ref, err := h.store.GetReviewKeysAndRating(reviewId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve review",
		})
		return
	} else if ref.UserId == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	} else if ref.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this review",
		})
		return
	}

	diff := req.Rating - ref.Rating
	resp, err := h.store.UpdateReview(req, reviewId, ref.PlaceId, diff)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Deletes a users review
// @Description	Deletes a users review given a valid access token
// @Tags        reviews
// @Accept      json
// @Produce     json
// @Param       id path int true "Review ID"
// @Success	    204
// @Router      /v1/reviews/{id} [delete]
// @Security    BearerAuth
func (h *Handler) deleteReview(c *gin.Context) {
    reviewId, err := strconv.Atoi(c.Param("id"))

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
        return
    }

    userId := c.MustGet("userId").(int)
	ref, err := h.store.GetReviewKeysAndRating(reviewId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retreive review",
		})
		return
	} else if ref.UserId == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	} else if ref.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this review",
		})
		return
	}

	err = h.store.DeleteReview(
		reviewId,
		ref.PlaceId,
		ref.Rating,
	)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(http.StatusNoContent, nil)
    }
}
