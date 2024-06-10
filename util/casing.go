package util

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	en     = language.English
	tcaser = cases.Title(en)
	ucaser = cases.Upper(en)
	scaser = cases.Lower(en)
)

func TitleCase(s string) string {
	s = strings.TrimSpace(s)
	return tcaser.String(s)
}

func LowerCase(s string) string {
	s = strings.TrimSpace(s)
	return scaser.String(s)
}

func SentenceCase(s string) string {
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")

	stopWords := strings.Join([]string{
		"a",
		"an",
		"and",
		"as",
		"at",
		"but",
		"by",
		"en",
		"for",
		"from",
		"how",
		"if",
		"in",
		"neither",
		"nor",
		"of",
		"on",
		"only",
		"onto",
		"out",
		"or",
		"per",
		"so",
		"than",
		"that",
		"the",
		"to",
		"until",
		"up",
		"upon",
		"v",
		"v.",
		"versus",
		"vs",
		"vs.",
		"via",
		"when",
		"with",
		"without",
		"yet",
	}, " ")

	for i, word := range words {
		// not capitalized if first letter of string
		if strings.Contains(stopWords, " "+word+" ") && word != string(word[0]) {
			words[i] = word
		} else {
			words[i] = tcaser.String(word)
		}
	}
	return strings.Join(words, " ")
}

func NameCase(s string) string {
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")

	var initials = regexp.MustCompile(`[a-zA-Z]\.[a-zA-Z]`)

	for i, word := range words {
		// handles initials
		if initials.MatchString(word) {
			words[i] = ucaser.String(word)
		} else {
			words[i] = tcaser.String(word)
		}
	}
	return strings.Join(words, " ")
}
