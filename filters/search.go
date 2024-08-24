package filters

type Search struct {
	Search string
	Filters
}

func (sf *Search) Empty() bool {
	return sf.Filters.Empty() &&
		sf.Search == ""
}
