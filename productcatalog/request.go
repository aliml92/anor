package productcatalog

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func isHXRequest(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	}

	return false
}

func extractProductID(slug string) (int64, error) {
	lastIndex := strings.LastIndex(slug, "-")
	if lastIndex == -1 {
		// assume slug is number string
		id, err := strconv.Atoi(slug)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", slug)
		}
		return int64(id), nil
	}

	if lastIndex == len(slug)-1 {
		return 0, fmt.Errorf("invalid slug format: %s", slug)
	}

	// Extract the ID part after the last "-"
	idStr := slug[lastIndex+1:]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ID to integer: %s", idStr)
	}

	return int64(id), nil
}

func extractCategoryID(slug string) (int32, error) {
	lastIndex := strings.LastIndex(slug, "-")
	if lastIndex == -1 {
		// assume slug is number string
		id, err := strconv.Atoi(slug)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", slug)
		}
		return int32(id), nil
	}

	if lastIndex == len(slug)-1 {
		return 0, fmt.Errorf("invalid slug format: %s", slug)
	}

	// Extract the ID part after the last "-"
	idStr := slug[lastIndex+1:]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ID to integer: %s", idStr)
	}

	return int32(id), nil
}
