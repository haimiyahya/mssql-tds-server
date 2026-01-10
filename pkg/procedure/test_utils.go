package procedure

import (
	"fmt"
	"testing"
)

func debugLog(t *testing.T, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	t.Log("DEBUG:", msg)
	fmt.Println("DEBUG:", msg)
}
