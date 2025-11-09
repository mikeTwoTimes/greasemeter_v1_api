package bookmarks

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
)

type Handler struct {
	store types.BookmarkStore
}

func NewHandler(store types.BookmarkStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.POST("/bookmarks/places/:id", h.createBookmark)
	auth.GET("/bookmarks", h.getBookmarksForUser)
	auth.GET("/bookmarks/places/:id", h.isPlaceBookmarked)
	auth.DELETE("/bookmarks/:id", h.deleteBookmark)
}

// @Summary	    Creates a bookmark for a place
// @Description	Creates a valid user's bookmark for a place
// @Tags        bookmarks
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    204
// @Router      /v1/bookmarks/places/{id} [post]
// @Security    BearerAuth
func (h *Handler) createBookmark(c *gin.Context) {
    placeId, err := strconv.Atoi(c.Param("id"))

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
        return
    }

    userId := c.MustGet("userId").(int)

	if err = h.store.CreateBookmark(userId, placeId); err != nil {
        c.JSON(utility.MapError(err))
    } else {
        c.JSON(http.StatusNoContent, nil)
    }
}

// @Summary	    Gets user's bookmarks
// @Description	Gets bookmarks for user given a valid access token
// @Tags        bookmarks
// @Accept      json
// @Produce     json
// @Success	    200 {object} []types.Bookmark
// @Router      /v1/bookmarks [get]
// @Security    BearerAuth
func (h *Handler) getBookmarksForUser(c *gin.Context) {
    userId := c.MustGet("userId").(int)
    resp, err := h.store.GetBookmarksForUser(userId)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Gets user's bookmark for place
// @Description	Gets a user's bookmark status of a place
// @Tags        bookmarks
// @Accept      json
// @Produce     json
// @Param       id path int true "Place ID"
// @Success	    200 boolean result
// @Router      /v1/bookmarks/places/{id} [get]
// @Security    BearerAuth
func (h *Handler) isPlaceBookmarked(c *gin.Context) {
	placeId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
		return
	}

	userId := c.MustGet("userId").(int)
	resp, err := h.store.IsPlaceBookmarked(userId, placeId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary	    Deletes a users bookmark
// @Description	Deletes a users bookmark given a valid access token
// @Tags        bookmarks
// @Accept      json
// @Produce     json
// @Param       id path int true "Bookmark ID"
// @Success	    204
// @Router      /v1/bookmarks/{id} [delete]
// @Security    BearerAuth
func (h *Handler) deleteBookmark(c *gin.Context) {
    bookmarkId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
        return
	}

	userId := c.MustGet("userId").(int)
	bookmarkUserId, err := h.store.GetUserFromBookmark(bookmarkId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if bookmarkUserId == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
	} else if userId != bookmarkUserId {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this bookmark",
		})
	} else if err = h.store.DeleteBookmark(bookmarkId); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(http.StatusNoContent, nil)
    }
}
