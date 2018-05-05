package user_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/gocapture/user"
)

func setupController(t *testing.T) (*bastion.Bastion, func()) {
	service, serviceTeardown := setupService(t)
	teardown := func() { serviceTeardown() }

	controller := user.Controller{
		Service: service,
		Render:  json.NewRender,
	}

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/users/", controller.Router())

	return app, teardown
}

func TestCreateValidUser(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:     "create user with only email",
			payload:  map[string]interface{}{"email": "test@localhost.com"},
			response: map[string]interface{}{"email": "test@localhost.com"},
		},
		{
			name: "create user",
			payload: map[string]interface{}{
				"email":    "test2@localhost.com",
				"password": "b4KeHAYy3u9v=ZQX",
			},
			response: map[string]interface{}{"email": "test2@localhost.com"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/users/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusCreated).
				JSON().Object().
				ContainsKey("email").ValueEqual("email", tc.response["email"]).
				ContainsKey("id").NotEmpty().
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty().
				NotContainsKey("password")
		})
	}
}

func TestCreateINValidUser(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "invalid email",
			payload: map[string]interface{}{"email": "test@localhost"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "invalid email",
			},
		},
		{
			name:    "missing email",
			payload: map[string]interface{}{"email": ""},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "email must not be blank",
			},
		},
		{
			name:    "missing email",
			payload: map[string]interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "email must not be blank",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/users/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestConflictEmail(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	payload := map[string]interface{}{"email": "test@localhost.com"}
	response := map[string]interface{}{
		"status":  409.0,
		"error":   "Conflict",
		"message": "email 'test@localhost.com' already exists",
	}

	e := bastion.Tester(t, app)
	e.POST("/users/").WithJSON(payload).Expect().Status(http.StatusCreated)

	e.POST("/users/").WithJSON(payload).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().Equal(response)
}
