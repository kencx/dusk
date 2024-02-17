package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	isbn = regexp.MustCompile(`(?:ISBN(?:-10|-13)?\s*:?\s*)?((?:(?:[\d]-?){9}|(?:[\d]-?){12})[\dxX])(?:[^-\d]|$)`)

	// correct format but invalid isbn digits
	ErrInvalidIsbn = errors.New("invalid isbn digits")
)

func IsbnCheck(value string) (bool, error) {
	if !isbn.MatchString(value) {
		return false, nil
	}

	value = strings.ReplaceAll(value, "-", "")
	for _, match := range isbn.FindAllStringSubmatch(value, -1) {
		if len(match) > 1 {

			switch i := match[1]; len(i) {
			case 10:
				if !isbn10Validate(i) {
					return false, ErrInvalidIsbn
				}
				return true, nil
			case 13:
				if !isbn13Validate(i) {
					return false, ErrInvalidIsbn
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func isbn10Validate(isbn string) bool {
	var sum int
	var mul = 10

	for i, c := range isbn {
		s := string(c)
		if i == 9 && s == "X" {
			s = "10"
		}

		d, err := strconv.Atoi(s)
		if err != nil {
			return false
		}

		sum += (d * mul)
		mul--
	}

	return sum%11 == 0
}

func isbn13Validate(isbn string) bool {
	var sum int

	for i, c := range isbn {
		var mul = 3
		if i%2 == 0 {
			mul = 1
		}

		d, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}

		sum += (d * mul)
	}

	return sum%10 == 0
}
