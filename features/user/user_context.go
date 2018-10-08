package user

import (
	"context"
	"errors"

	"github.com/ifreddyrondon/capture/features"
	"gopkg.in/src-d/go-kallax.v1"
)

type ctxKey string

const userKey ctxKey = "user"

var (
	errMissingUser    = errors.New("user not found in context")
	errWrongUserValue = errors.New("user value set incorrectly in context")
)

// WithUser will return a new context with the user value added to it.
func WithUser(ctx context.Context, user *features.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetFromContext returns the user assigned to the context, or error if there
// is any error or there isn't a user.
func GetFromContext(ctx context.Context) (*features.User, error) {
	tmp := ctx.Value(userKey)
	if tmp == nil {
		return nil, errMissingUser
	}
	user, ok := tmp.(*features.User)
	if !ok {
		return nil, errWrongUserValue
	}
	return user, nil
}

// GetUserID will return the user ID assigned to the context, or error if there
// is any error or there isn't a user.
func GetUserID(ctx context.Context) (kallax.ULID, error) {
	u, err := GetFromContext(ctx)
	if err != nil {
		return kallax.ULID{}, err
	}
	return u.ID, nil
}
