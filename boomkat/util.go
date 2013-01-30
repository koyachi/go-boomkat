package boomkat

import (
	"strings"
)

func GenresFromString(text string) []string {
	genres := strings.Split(text, "/")
	for i, genre := range genres {
		genres[i] = strings.TrimSpace(genre)
	}
	return genres
}
