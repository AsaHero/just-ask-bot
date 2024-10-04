package utility

import "regexp"

func IsValidEmail(email string) bool {
	// Regular expression for validating email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the email matches the regex pattern
	return regex.MatchString(email)
}
