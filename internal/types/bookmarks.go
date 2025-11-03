package types

type BookmarkStore interface {
	InsertBookmark(userId int, placeId int) error
	GetUserByBookmark(bookmarkId int) (int, error)
	GetBookmarksByUser(userId int) ([]Bookmark, error)
	DeleteBookmark(bookmarkId int) error
}

type Bookmark struct {
	Id      int     `json:"id"`
	PlaceId int     `json:"placeId"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
}
