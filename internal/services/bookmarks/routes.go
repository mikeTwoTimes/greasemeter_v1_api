package bookmarks

import (
	"greasemeter_v1_api/internal/types"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.BookmarkStore
}

func NewHandler(store types.BookmarkStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/bookmarks")
	{
		v1.POST("/places/:id", h.createBookmark)
		v1.GET("/", h.getBookmarksByUser)
		v1.DELETE("/:id", h.deleteBookmark)
	}
}

func (h *Handler) createBookmark(c *gin.Context) {

}

func (h *Handler) getBookmarksFromUser(c *gin.Context) {

}

func (h *Handler) deleteBookmark(c *gin.Context) {

}
