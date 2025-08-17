package state

type Tag struct {
	Base
}

// func NewTagState(r *http.Request) (*Base, validator.ErrMap) {
// 	view := r.URL.Query().Get("view")
//
// 	filters := initTagFilters(r)
// 	if errMap := validator.Validate(filters.Base); errMap != nil {
// 		return nil, errMap
// 	}
//
// 	return NewBase(view, filters), nil
// }
