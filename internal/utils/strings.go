package utils

import (
	"strings"
)

func CheckTranslation(translations []string, guess string) bool {
	clearedGuess := strings.ToLower(strings.TrimSpace(guess))

	for _, translation := range translations {
		clearedTranslation := strings.ToLower(translation)

		if clearedTranslation == clearedGuess {
			return true
		}
	}

	return false
}
