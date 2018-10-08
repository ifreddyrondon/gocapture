package decoder_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/user/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestDecodePostUserOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name        string
		emailResult string
		body        string
	}{
		{
			name:        "decode user without password",
			emailResult: "test@example.com",
			body:        `{"email": "test@example.com"}`,
		},
		{
			name:        "decode user with password",
			emailResult: "test@example.com",
			body:        `{"email":"test@example.com","password":"1234"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var u decoder.PostUser
			err := decoder.Decode(r, &u)
			assert.Nil(t, err)
			assert.Equal(t, tc.emailResult, *u.Email)
		})
	}
}

func TestDecodePostUserError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "decode user's email missing",
			body: `{}`,
			err:  "email must not be blank",
		},
		{
			name: "decode user's email missing",
			body: `{"email": "test@"}`,
			err:  "invalid email",
		},
		{
			name: "decode user's password too short",
			body: `{"email":"test@example.com","password":"1"}`,
			err:  "password must have at least four characters",
		},
		{
			name: "invalid user payload",
			body: `.`,
			err:  "cannot unmarshal json into valid user",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var u decoder.PostUser
			err := decoder.Decode(r, &u)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestUserPostUserOK(t *testing.T) {
	email, pass := "test@example.com", "1234"
	t.Parallel()
	tt := []struct {
		name     string
		postUser decoder.PostUser
	}{
		{
			name:     "get user from postUser without password",
			postUser: decoder.PostUser{Email: &email},
		},
		{
			name:     "get user from postUser with password",
			postUser: decoder.PostUser{Email: &email, Password: &pass},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var u features.User
			err := decoder.User(tc.postUser, &u)
			assert.Nil(t, err)
			// test user fields filled with not default values
			assert.NotEqual(t, kallax.ULID{}, u.ID)
			assert.NotEqual(t, time.Time{}, u.CreatedAt)
			assert.NotEqual(t, time.Time{}, u.UpdatedAt)
		})
	}
}