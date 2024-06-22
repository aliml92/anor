package product

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	minSearchQueryLength = 1
	maxSearchQueryLength = 256
)

var (
	ErrInvalidSearchQueryLength = errors.New("search query must be between 1 and 256 characters")
	ErrInvalidSearchQueryChars  = errors.New("search query contains invalid characters")
)

func validateCategoryHandle(handle string) error {
	// Check if the handle is directly a number string
	if _, err := strconv.Atoi(handle); err == nil {
		// If no error, the handle is a valid number string
		return nil
	}

	// If the handle is not directly a number, check for the '-' separator
	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 || lastIndex == len(handle)-1 {
		// If '-' is not found or is the last character, handle format is invalid
		return fmt.Errorf("invalid handle format: %s", handle)
	}

	// Extract the ID part after the last "-"
	idStr := handle[lastIndex+1:]

	// Validate the ID part is a number string
	if _, err := strconv.Atoi(idStr); err != nil {
		return fmt.Errorf("invalid ID format in handle: %s", handle)
	}

	// If no issues, the handle format is valid
	return nil
}

func validateSearchQuery(q string) error {
	if len(q) < minSearchQueryLength || len(q) > maxSearchQueryLength {
		return ErrInvalidSearchQueryLength
	}

	// This pattern allows letters (including those from various scripts/languages),
	// numbers, spaces, and some punctuation. It uses the Unicode character property classes.
	pattern := `^[\p{L}\p{N}\s.,-]+$`
	matched, err := regexp.MatchString(pattern, q)
	if err != nil {
		return err
	}

	if !matched {
		return ErrInvalidSearchQueryChars
	}

	return nil
}
