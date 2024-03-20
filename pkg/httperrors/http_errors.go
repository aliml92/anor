package httperrors

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

const (
	PACKAGENAME = "github.com/aliml92/anor"
)

func ClientError(logger *slog.Logger, w http.ResponseWriter, err error, statusCode int) {
	_, file, no, _ := runtime.Caller(2)
	logger.LogAttrs(
		context.TODO(),
		slog.LevelError,
		"client error",
		slog.String("file", extractFilePath(file)),
		slog.String("line", strconv.Itoa(no)),
		slog.String("status", strconv.Itoa(statusCode)),
		slog.String("error", err.Error()),
	)

	http.Error(w, err.Error(), statusCode)
}

func ServerInternalError(logger *slog.Logger, w http.ResponseWriter, err error) {
	_, file, no, _ := runtime.Caller(2)
	logger.LogAttrs(
		context.TODO(),
		slog.LevelError,
		"server error",
		slog.String("file", file),
		slog.String("line", strconv.Itoa(no)),
		slog.String("status", strconv.Itoa(http.StatusInternalServerError)),
		slog.String("error", err.Error()),
	)

	http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
}

func LogClientError(logger *slog.Logger, err error) {
	_, file, no, _ := runtime.Caller(2)
	logger.LogAttrs(
		context.TODO(),
		slog.LevelError,
		"client error",
		slog.String("file", extractFilePath(file)),
		slog.String("line", strconv.Itoa(no)),
		slog.String("error", err.Error()),
	)
}

func extractFilePath(filePath string) string {
	index := strings.Index(filePath, PACKAGENAME)
	if index == -1 {
		return filePath
	}
	trimmedPath := filePath[index+len(PACKAGENAME)+1:]
	return trimmedPath
}
