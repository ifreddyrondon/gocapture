package collection_test

import (
	"sync"
	"testing"

	"github.com/ifreddyrondon/gocapture/collection"
	"github.com/ifreddyrondon/gocapture/database"
	"github.com/jinzhu/gorm"
)

var once sync.Once
var db *gorm.DB

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setupRepository(t *testing.T) (collection.Repository, func()) {
	repo := collection.NewPGRepository(getDB().Table("collections"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	return repo, teardown
}
