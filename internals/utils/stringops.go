package utils

import (
	"strings"
)

func SlugifyBase(in string) string {
	ot := strings.ToLower(in)
	rep := strings.NewReplacer(" ", "-", "+", "-", "_", "-", ".", "-", "/", "-")
	return rep.Replace(ot)
}

func SetIds(separator string, opts ...string) string {
	ids := ""
	for _, val := range opts {
		if ids == "" {
			ids = val
		} else {
			ids = ids + separator + val
		}
	}
	return ids
}
