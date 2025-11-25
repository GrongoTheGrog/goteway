package utils

import (
	"bytes"
	"io"
	"net/http"
)

func ErrorResponse(message string, status int) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer([]byte(message))),
		StatusCode: status,
		Status:     string(rune(status)),
	}
}
