package block

import "fmt"

func prefix(prefix, value string) string {
	return fmt.Sprintf("%s_%s", prefix, value)
}
