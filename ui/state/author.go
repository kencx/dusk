package state

type Author struct {
	Base
}

// func NewAuthorState(r *http.Request) (*Base, validator.ErrMap) {
// 	view := r.URL.Query().Get("view")
//
// 	filters := initAuthorFilters(r)
// 	if errMap := validator.Validate(filters.Base); errMap != nil {
// 		return nil, errMap
// 	}
//
// 	return NewBase(view, filters), nil
// }
