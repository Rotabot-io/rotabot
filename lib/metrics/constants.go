package metrics

import "fmt"

func Endpoint(verb, path string) string {
	return fmt.Sprintf("%s:%s", verb, path)
}
