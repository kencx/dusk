package validator

import "regexp"

type Validator struct {
	Errors map[string]string `json:"errors"`
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func Matches(value string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(value)
}

func In[T comparable](item T, sl []T) bool {
	for _, elem := range sl {
		if item == elem {
			return true
		}
	}
	return false
}
