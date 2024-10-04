package utility

import (
	"fmt"
	"strings"
)

// SplitIntoParagraphs splits the given text into a slice of paragraphs based on empty lines.
func SplitIntoParagraphs(text string) []string {
	// Normalize line endings to Unix style if unsure about the input format.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.Split(text, "\n\n")
}

func Pluralize(number int, singular string, plural2_4 string, plural5_0 string) string {
	// Get the last two digits
	lastTwoDigits := number % 100
	// Get the last digit
	lastDigit := number % 10

	// Determine the correct form
	if lastTwoDigits >= 11 && lastTwoDigits <= 14 {
		return plural5_0
	}

	switch lastDigit {
	case 1:
		return singular
	case 2, 3, 4:
		return plural2_4
	default:
		return plural5_0
	}
}

func DashFormat(strs []string) string {
	var result string
	for _, str := range strs {
		result += fmt.Sprintf("- %s\n", str)
	}

	return result
}
