package anor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/samber/oops"
	"log/slog"
	"net/http"
	"strings"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUserNotFound = errors.New("user not found")

	ErrPaymentNotFound               = errors.New("payment not found")
	ErrUserExists                    = errors.New("user exists")
	ErrProductNotFound               = errors.New("product not found")
	ErrProductPricingNotFound        = errors.New("product pricing not found")
	ErrCartNotFound                  = errors.New("cart not found")
	ErrCartItemNotFound              = errors.New("cart item not found")
	ErrProductVariantNotFound        = errors.New("product variant not found")
	ErrProductVariantPricingNotFound = errors.New("product variant pricing not found")
	ErrAtLeastOneOptionRequired      = errors.New("at least one option required")
	ErrAddressNotFound               = errors.New("address not found")
)

func ClientError(logger *slog.Logger, w http.ResponseWriter, err error, statusCode int) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
	http.Error(w, err.Error(), statusCode)
}

func ServerInternalError(logger *slog.Logger, w http.ResponseWriter, err error) {
	// TODO: slog not render \n in stacktrace error
	//logger.Error(
	//	err.Error(),
	//	slog.Any("error", err),
	//)
	var oopsErr oops.OopsError
	if ok := errors.As(err, &oopsErr); ok {
		fmt.Println(colorizeStacktrace(oopsErr.Stacktrace()))
	} else {
		fmt.Println(err)
	}

	http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
}

func LogError(logger *slog.Logger, err error) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
}

func colorizeStacktrace(stacktrace string) string {
	lines := strings.Split(stacktrace, "\n")
	colorizedLines := make([]string, 0, len(lines))

	for i, line := range lines {
		if i == 0 {
			// Colorize the main error message
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				colorizedLine := fmt.Sprintf("%s: %s", parts[0], color.RedString(parts[1]))
				colorizedLines = append(colorizedLines, colorizedLine)
			} else {
				colorizedLines = append(colorizedLines, line)
			}
		} else if strings.Contains(line, "at ") {
			// Colorize the file and function information
			parts := strings.SplitN(line, "at ", 2)
			if len(parts) == 2 {
				fileAndFunc := strings.SplitN(parts[1], " ", 2)
				if len(fileAndFunc) == 2 {
					lastSlash := strings.LastIndex(fileAndFunc[0], "/")
					path := fileAndFunc[0][:lastSlash+1]
					fileAndLine := fileAndFunc[0][lastSlash+1:]

					colorizedLine := fmt.Sprintf("%s at %s%s %s",
						parts[0],
						path,
						color.CyanString(fileAndLine),
						color.GreenString(fileAndFunc[1]))
					colorizedLines = append(colorizedLines, colorizedLine)
				} else {
					colorizedLines = append(colorizedLines, line)
				}
			} else {
				colorizedLines = append(colorizedLines, line)
			}
		} else {
			colorizedLines = append(colorizedLines, line)
		}
	}

	return strings.Join(colorizedLines, "\n")
}

// JSONErrorResponse represents a standardized JSON error response
type JSONErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code,omitempty"`
}

// NewJSONErrorResponse creates a new JSONErrorResponse
func NewJSONErrorResponse(err error, statusCode int) JSONErrorResponse {
	return JSONErrorResponse{
		Error:      err.Error(),
		StatusCode: statusCode,
	}
}

// JSONServerInternalError writes a JSON error response to the http.ResponseWriter
func JSONServerInternalError(logger *slog.Logger, w http.ResponseWriter, err error) {
	// TODO: slog not render \n in stacktrace error
	//logger.Error(
	//	err.Error(),
	//	slog.Any("error", err),
	//)
	var oopsErr oops.OopsError
	if ok := errors.As(err, &oopsErr); ok {
		fmt.Println(colorizeStacktrace(oopsErr.Stacktrace()))
	} else {
		fmt.Println(err)
	}

	errResp := NewJSONErrorResponse(err, http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errResp)
}

func JSONClientError(logger *slog.Logger, w http.ResponseWriter, err error, statusCode int) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
	http.Error(w, err.Error(), statusCode)

	errResp := NewJSONErrorResponse(err, statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResp)
}
