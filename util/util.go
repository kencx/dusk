package util

import (
	"strings"

	"github.com/kencx/dusk/null"
)

// trim \n, \t
func TrimMultiLine(s string) string {
	q := strings.Split(s, "\n")
	for i, qi := range q {
		qi = strings.TrimSpace(qi)
		q[i] = qi
	}
	return strings.Join(q, " ")
}

func PrintDateMonthYear(date null.Time) string {
	return PrintDateFormat(date, "Jan 2006")
}

func PrintDateFull(date null.Time) string {
	return PrintDateFormat(date, "02 Jan 2006")
}

func PrintDateFormat(date null.Time, format string) string {
	return date.ValueOrZero().Format(format)
}
