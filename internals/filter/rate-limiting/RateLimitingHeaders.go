package rate_limiting

import (
	"net/http"
	"strconv"
)

func writeRateLimitingHeaders(response *http.Response, limit, remaining int) {
	response.Header.Set("X-RateLimit-Limit", strconv.Itoa(limit))
	response.Header.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
}
