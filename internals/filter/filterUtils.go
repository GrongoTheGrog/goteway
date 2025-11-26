package filter

import "fmt"

func RequestPrefix(requestId string) string {
	return fmt.Sprintf("REQUEST %s:", requestId)
}
