package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/pkg/errors"
)

var (
	// RepoCtxKey is the context.Context key to store the Repo for a request.
	RepoCtxKey = &contextKey{"Repository"}
)
var (
	errMissingCtxRepo = errors.New("repo not found in context")
	errWrongRepoValue = errors.New("repo value set incorrectly in context")
	errMissingRepo    = errors.New("not found repository")
	errInvalidRepoID  = errors.New("invalid repository id")
)

func withRepo(ctx context.Context, repo *pkg.Repository) context.Context {
	return context.WithValue(ctx, RepoCtxKey, repo)
}

// GetRepo returns the repo assigned to the context, or error if there
// is any error or there isn't a repo.
func GetRepo(ctx context.Context) (*pkg.Repository, error) {
	tmp := ctx.Value(RepoCtxKey)
	if tmp == nil {
		return nil, errMissingCtxRepo
	}
	repo, ok := tmp.(*pkg.Repository)
	if !ok {
		return nil, errWrongRepoValue
	}
	return repo, nil
}

func RepoCtx(service getting.Service) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			repoID := chi.URLParam(r, "id")
			u, err := GetUser(r.Context())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				json.InternalServerError(w, err)
				return
			}

			repo, err := service.GetRepo(repoID, u)
			if err != nil {
				if isInvalidErr(err) {
					json.BadRequest(w, errInvalidRepoID)
					return
				}
				if isNotFound(err) {
					json.NotFound(w, errMissingRepo)
					return
				}
				if isNotAuthorized(err) {
					httpErr := render.HTTPError{
						Status:  http.StatusForbidden,
						Error:   http.StatusText(http.StatusForbidden),
						Message: "not authorized to see this repository",
					}
					json.Response(w, http.StatusForbidden, httpErr)
					return
				}
				fmt.Fprintln(os.Stderr, err)
				json.InternalServerError(w, err)
				return
			}

			ctx := withRepo(r.Context(), repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
