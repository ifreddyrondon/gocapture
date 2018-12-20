package middleware

import "github.com/pkg/errors"

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "capture/middleware context value " + k.name
}

type invalidErr interface {
	IsInvalid() bool
}

func isInvalidErr(err error) bool {
	if e, ok := errors.Cause(err).(invalidErr); ok {
		return e.IsInvalid()
	}
	return false
}

type authorizationErr interface {
	// IsNotAllowed returns true when the req is not allowed.
	IsNotAuthorized() bool
}

func isNotAuthorized(err error) bool {
	if e, ok := errors.Cause(err).(authorizationErr); ok {
		return e.IsNotAuthorized()
	}
	return false
}

type notFoundErr interface {
	// NotFound returns true when a resource is not found.
	NotFound() bool
}

func isNotFound(err error) bool {
	if e, ok := errors.Cause(err).(notFoundErr); ok {
		return e.NotFound()
	}
	return false
}
