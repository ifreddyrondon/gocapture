package middleware_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var tempUser = pkg.User{Email: "test@example.com", ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}

func authenticatedMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserCtxKey, &tempUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func notUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func setupRepoCtx(service getting.Service, auth func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(auth)
		r.Use(middleware.RepoCtx(service))
		r.Get("/{id}", handler)
		r.Post("/{id}", handler)
	})
	return app
}

type mockGettingService struct {
	repo *pkg.Repository
	err  error
}

func (m *mockGettingService) GetRepo(string, *pkg.User) (*pkg.Repository, error) {
	return m.repo, m.err
}

func TestRepoCtxSuccess(t *testing.T) {
	t.Parallel()

	app := setupRepoCtx(&mockGettingService{}, authenticatedMiddle)
	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusOK)
}

func TestRepoCtxFailInternalErrorGettingUser(t *testing.T) {
	t.Parallel()
	s := &mockGettingService{}
	app := setupRepoCtx(s, notUserMiddleware)
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailBadRequestGettingRepoByInvalidErr(t *testing.T) {
	t.Parallel()
	s := &mockGettingService{err: invalidErr("test")}
	app := setupRepoCtx(s, authenticatedMiddle)
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid repository id",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailNotFoundGettingRepo(t *testing.T) {
	t.Parallel()
	s := &mockGettingService{err: notFound("test")}
	app := setupRepoCtx(s, authenticatedMiddle)
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "not found repository",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailForbiddenGettingRepo(t *testing.T) {
	t.Parallel()
	s := &mockGettingService{err: notAllowedErr("test")}
	app := setupRepoCtx(s, authenticatedMiddle)
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "not authorized to see this repository",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailInternalServerErrGettingRepo(t *testing.T) {
	t.Parallel()
	s := &mockGettingService{err: errors.New("test")}
	app := setupRepoCtx(s, authenticatedMiddle)
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestContextGetRepoOK(t *testing.T) {
	ctx := context.Background()
	repo := pkg.Repository{Name: "test"}
	ctx = context.WithValue(ctx, middleware.RepoCtxKey, &repo)

	r, err := middleware.GetRepo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, r.Name, r.Name)
}

func TestContextGetRepoMissingRepo(t *testing.T) {
	ctx := context.Background()
	_, err := middleware.GetRepo(ctx)
	assert.EqualError(t, err, "repo not found in context")
}

func TestContextGetRepoIDMissingRepo(t *testing.T) {
	ctx := context.Background()

	_, err := middleware.GetRepo(ctx)
	assert.EqualError(t, err, "repo not found in context")
}
