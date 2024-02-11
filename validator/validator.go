package validator

import (
	"fmt"
	"regexp"
	"strings"
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

	var msg strings.Builder
	for k, v := range e {
		msg.WriteString(fmt.Sprintf("%s=%s, ", k, v))

	}
	return msg.String()
}

func (e ErrMap) Check(ok bool, key, message string) {
	if !ok {
		e.Add(key, message)
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
