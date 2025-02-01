package validator

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
)

type Validator interface {
	Valid() ErrMap
}

type ErrMap map[string]string

func New() ErrMap {
	return ErrMap(make(map[string]string))
}

func (e ErrMap) Error() string {
	if len(e) == 0 {
		return ""
	}

	// returns raw json string of error map
	mJson, err := json.Marshal(e)
	if err != nil {
		slog.Error("[Validator] Failed to marshal error map", slog.Any("err", err))
		return ""
	}

	return string(mJson)
}

func (e ErrMap) Check(ok bool, key, message string) {
	if !ok {
		e.Add(key, message)
	}
}

func (e ErrMap) EitherOr(either, or bool, key1, key2, message string) {
	if !either && !or {
		e.Add(fmt.Sprintf("%s or %s", key1, key2), message)
	}
}

func (e ErrMap) Add(key, message string) {
	if _, exists := e[key]; !exists {
		e[key] = message
	}
}

func Validate(v Validator) ErrMap {
	if err := v.Valid(); len(err) > 0 {
		return err
	}
	return nil
}

func Matches(value string, r *regexp.Regexp) bool {
	return r.MatchString(value)
}

func In[T comparable](item T, sl []T) bool {
	for _, elem := range sl {
		if item == elem {
			return true
		}
	}
	return false
}
