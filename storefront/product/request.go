package product

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

func extractProductID(handle string) (int64, error) {
	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 {
		// assume handle is number string
		id, err := strconv.Atoi(handle)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", handle)
		}
		return int64(id), nil
	}

	if lastIndex == len(handle)-1 {
		return 0, fmt.Errorf("invalid handle format: %s", handle)
	}

	// Extract the ID part after the last "-"
	idStr := handle[lastIndex+1:]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ID to integer: %s", idStr)
	}

	return int64(id), nil
}

func extractCategoryID(handle string) (int32, error) {
	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 {
		// assume handle is number string
		id, err := strconv.Atoi(handle)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", handle)
		}
		return int32(id), nil
	}

	if lastIndex == len(handle)-1 {
		return 0, fmt.Errorf("invalid handle format: %s", handle)
	}

	// Extract the ID part after the last "-"
	idStr := handle[lastIndex+1:]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ID to integer: %s", idStr)
	}

	return int32(id), nil
}

func extractID(handle string) int64 {
	parts := strings.Split(handle, "-")
	if len(parts) > 1 {
		idStr := parts[len(parts)-1]
		id, _ := strconv.ParseInt(idStr, 10, 64) // Ignore the error
		return id
	}

	id, _ := strconv.ParseInt(handle, 10, 64) // Ignore the error
	return id
}
