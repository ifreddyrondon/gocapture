package repository

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Controller handler the repository routes
type Controller struct {
	service Service
	render  render.Render
}

// NewController returns a new Controller
func NewController(service Service, render render.Render) *Controller {
	return &Controller{service: service, render: render}
}

// Router creates a REST router for the user resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Post("/", c.create)
	return r
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		_ = c.render(w).BadRequest(err)
		return
	}

	if err := c.service.Save(&repo); err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}

	_ = c.render(w).Created(repo)
}