package filters

type Search struct {
	Search string
	Base
}

func (sf *Search) Empty() bool {
	return sf.Base.Empty() &&
		sf.Search == ""
}
