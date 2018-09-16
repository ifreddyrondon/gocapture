package auth

import (
	"net/http"

	"github.com/ifreddyrondon/capture/features/auth/jwt"

	"github.com/ifreddyrondon/capture/features/user"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type tokenJSON struct {
	Token string `json:"token,omitempty"`
}

// Controller handler the auth routes
type Controller struct {
	authenticate func(http.Handler) http.Handler
	render       render.APIRenderer
	service      *jwt.Service
}

// NewController returns a new Controller
func NewController(authenticate func(http.Handler) http.Handler, service *jwt.Service) *Controller {
	return &Controller{
		authenticate: authenticate,
		service:      service,
		render:       render.NewJSON(),
	}
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/token-auth", func(r chi.Router) {
		r.Use(c.authenticate)
		r.Post("/", c.login)
	})
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	u, err := user.GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	token, err := c.service.GenerateToken(u.ID.String())
	if err != nil {
		c.render.InternalServerError(w, err)
	}

	c.render.Send(w, tokenJSON{Token: token})
}
