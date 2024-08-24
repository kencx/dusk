package filters

type Book struct {
	Title  string
	Author string
	Tag    string
	Series string
	Search
}

func (bf *Book) Empty() bool {
	return bf.Search.Empty() &&
		bf.Title == "" &&
		bf.Author == "" &&
		bf.Tag == "" &&
		bf.Series == ""
}
