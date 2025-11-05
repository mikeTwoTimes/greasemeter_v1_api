package types

type BookmarkStore interface {
	CreateBookmark(userId, placeId int) error
	GetUserFromBookmark(bookmarkId int) (int, error)
	GetBookmarksForUser(userId int) ([]Bookmark, error)
	IsPlaceBookmarked(userId, placeId int) (bool, error)
	DeleteBookmark(bookmarkId int) error
}

type Bookmark struct {
	Id      int     `json:"id"`
	PlaceId int     `json:"placeId"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
}
