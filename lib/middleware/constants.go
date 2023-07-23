package middleware

type key string

// RequestIdKey is the key used to store the request ID in the context or header.
const RequestIdKey key = "X-Request-Id"
