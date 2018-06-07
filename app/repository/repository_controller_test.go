package repository_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/app/auth/authorization"

	"github.com/ifreddyrondon/capture/app/repository"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

type MockService struct{}

func (r *MockService) Save(c *repository.Repository) error {
	return errors.New("test")
}

type mockStrategySuccess struct{}

func (m *mockStrategySuccess) IsAuthorizedREQ(r *http.Request) (string, error) {
	return "test", nil
}
func (m *mockStrategySuccess) IsNotAuthorizedErr(err error) bool {
	return false
}

type mockStrategyFail struct{}

func (m *mockStrategyFail) IsAuthorizedREQ(r *http.Request) (string, error) {
	return "", errors.New("you don’t have permission to access this resource")
}
func (m *mockStrategyFail) IsNotAuthorizedErr(err error) bool {
	return true
}

func setupControllerMockService(strategy authorization.Strategy) *bastion.Bastion {
	service := &MockService{}
	midl := authorization.NewAuthorization(strategy, json.NewRender)
	controller := repository.NewController(service, json.NewRender, midl)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/repository/", controller.Router())
	return app
}

func setupController(t *testing.T, strategy authorization.Strategy) (*bastion.Bastion, func()) {
	service, teardown := setupService(t)
	midl := authorization.NewAuthorization(strategy, json.NewRender)
	controller := repository.NewController(service, json.NewRender, midl)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/repository/", controller.Router())

	return app, teardown
}

func TestCreateRepositorySuccess(t *testing.T) {
	app, teardown := setupController(t, &mockStrategySuccess{})
	defer teardown()

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/repository/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("name").ValueEqual("name", payload["name"]).
		ContainsKey("shared").ValueEqual("shared", false).
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestCreateRepositoryFail(t *testing.T) {
	app, teardown := setupController(t, &mockStrategySuccess{})
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "no data",
			payload: map[string]interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "name must not be blank",
			},
		},
		{
			name:    "empty name",
			payload: map[string]interface{}{"name": ""},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "name must not be blank",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/repository/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestCreateRepositorySaveFail(t *testing.T) {
	t.Parallel()

	app := setupControllerMockService(&mockStrategySuccess{})

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/repository/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()
}

func TestCreateRepositoryNotAuthorized(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "you don’t have permission to access this resource",
	}

	app := setupControllerMockService(&mockStrategyFail{})
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/repository/").
		WithJSON(payload).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}
