package dusk

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Authors []*Author
