package middleware

import (
	"net/http"
	"path"
)

func requestMethod(r *http.Request) string {
	return r.Method
}

func requestPath(r *http.Request) string {
	return path.Clean(r.URL.EscapedPath())
}
