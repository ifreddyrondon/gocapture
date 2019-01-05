package getting_test

import (
	"fmt"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type authorizationErr interface{ IsNotAuthorized() bool }

type mockRepoStore struct {
	repo *domain.Repository
	err  error
}

func (m *mockRepoStore) Get(kallax.ULID) (*domain.Repository, error) { return m.repo, m.err }

func TestServiceGetRepoOKWhenUserOwner(t *testing.T) {
	t.Parallel()

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	repoID := kallax.NewULID()
	store := &mockRepoStore{repo: &domain.Repository{ID: repoID, Name: "test1", UserID: userID}}
	s := getting.NewRepoService(store)

	u := &domain.User{ID: userID}
	repo, err := s.Get(repoID, u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoOKWhenPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockRepoStore{repo: &domain.Repository{ID: repoID, Name: "test1", Visibility: domain.Public}}
	s := getting.NewRepoService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	repo, err := s.Get(repoID, u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoErrWhenNoOwnerAndNoPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockRepoStore{repo: &domain.Repository{ID: repoID, Name: "test1", Visibility: domain.Private}}
	s := getting.NewRepoService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	_, err = s.Get(repoID, u)
	assert.EqualError(t, err, fmt.Sprintf("user %v not authorized to get repo %v", userID, repoID))
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}

func TestServiceGetRepoErrGettingRepoFromStorage(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{err: errors.New("test")}
	s := getting.NewRepoService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	_, err = s.Get(kallax.NewULID(), u)
	assert.EqualError(t, err, "could not get repo: test")
}