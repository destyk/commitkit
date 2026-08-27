package lint

import "fmt"

func formatString(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
