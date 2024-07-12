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
	if _, err := strconv.Atoi(handle); err == nil {
		// If no error, the handle is a valid number string
		return nil
	}

	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 || lastIndex == len(handle)-1 {
		// If '-' is not found or is the last character, handle format is invalid
		return fmt.Errorf("invalid handle format: %s", handle)
	}

	idStr := handle[lastIndex+1:]
	if _, err := strconv.Atoi(idStr); err != nil {
		return fmt.Errorf("invalid ID format in handle: %s", handle)
	}

	return nil
}

func validateSearchQuery(q string) error {
	if len(q) < minSearchQueryLength || len(q) > maxSearchQueryLength {
		return ErrInvalidSearchQueryLength
	}

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
