package dusk

type Tag struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Books       Books  `json:"books,omitempty"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type Tags []*Tag
