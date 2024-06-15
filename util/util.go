package util

import "strings"

// trim \n, \t
func TrimMultiLine(s string) string {
	q := strings.Split(s, "\n")
	for i, qi := range q {
		qi = strings.TrimSpace(qi)
		q[i] = qi
	}
	return strings.Join(q, " ")
}
