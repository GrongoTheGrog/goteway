package filter

import "fmt"

func requestPrefix(requestId string) string {
	return fmt.Sprintf("REQUEST %s:", requestId)
}
