package utils

import (
	"strings"
	"unicode"
)

// ContainsAllWords checks if 'hay' contains every word in the 'needles' slice.
// It performs a case-insensitive, whole-word match.
// containsAllWords checks if 'text' contains every word in the 'words' slice.
// It performs a case-insensitive check.
func ContainsAllWords(hay string, needles []string) bool {
	// Convert all text to lowercase for case-insensitive comparison
	hay = strings.ToLower(hay)

	// Split the text into individual words based on whitespace
	// We filter out empty strings that might result from multiple spaces
	textWords := []string{}
	currentWord := ""
	for _, r := range hay {
		if unicode.IsSpace(r) {
			if currentWord != "" {
				textWords = append(textWords, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(r)
		}
	}
	if currentWord != "" {
		textWords = append(textWords, currentWord)
	}

	// Check if every word in the target array exists in the text words
	for _, targetWord := range needles {
		targetWord = strings.ToLower(targetWord)

		// Ensure we are matching the whole word (not a substring of a longer word)
		// We check if the word is in the list AND doesn't have a letter before or after it in the text
		found := false
		for _, tWord := range textWords {
			if tWord == targetWord {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
