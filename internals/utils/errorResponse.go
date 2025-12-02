package utils

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

func ErrorResponse(message string, status int) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer([]byte(message))),
		StatusCode: status,
		Status:     strconv.Itoa(status),
		Header:     http.Header{},
	}
}
