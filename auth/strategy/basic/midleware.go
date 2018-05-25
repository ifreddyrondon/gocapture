package basic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/user"
)

var (
	errInvalidCredentials = errors.New("invalid email or password")
	errInvalidPassword    = errors.New("invalid password")
)

// Strategy is a basic authentication method that uses email and password to authenticate
type Strategy struct {
	render.Render
	user.GetterService
	UserKey fmt.Stringer
}

// Authenticate for basic (username/password) authentication.
func (s *Strategy) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cre Crendentials
		if err := json.NewDecoder(r.Body).Decode(&cre); err != nil {
			_ = s.Render(w).BadRequest(err)
			return
		}

		u, err := s.validate(&cre)
		if err != nil {
			if err == errInvalidPassword || err == user.ErrNotFound {
				httpErr := bastionJSON.HTTPError{
					Status:  http.StatusUnauthorized,
					Errors:  http.StatusText(http.StatusUnauthorized),
					Message: errInvalidCredentials.Error(),
				}
				_ = s.Render(w).Response(http.StatusUnauthorized, httpErr)
				return
			}

			_ = s.Render(w).InternalServerError(err)
			return
		}
		ctx := context.WithValue(r.Context(), s.UserKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Strategy) validate(cre *Crendentials) (*user.User, error) {
	u, err := s.GetterService.Get(cre.Email)
	if err != nil {
		return nil, err
	}
	if !u.CheckPassword(cre.Password) {
		return nil, errInvalidPassword
	}

	return u, nil
}