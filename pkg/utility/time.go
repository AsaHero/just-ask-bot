package utility

import (
	"strings"
	"time"
)

// StartOfDate formats date time in 00:00:00
func StartOfDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// IsValidDateLayout checks if the date string is in the given format.
func IsValidDateLayout(dateStr string, format string) bool {
	_, err := time.Parse(format, dateStr)
	if err != nil && strings.Contains(err.Error(), "cannot parse") {
		return false // Format error
	}
	return true // Format is correct
}

// CheckDateValues checks if the date values in the string are valid.
func IsValidDateValue(dateStr string, format string) bool {
	_, err := time.Parse(format, dateStr)
	return err == nil // Returns true if there's no error, meaning the date values are valid
}
