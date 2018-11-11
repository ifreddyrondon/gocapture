package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/config"
	"github.com/sarulabs/di"
)

func router(resources di.Container) http.Handler {
	r := chi.NewRouter()

	authRoutes := resources.Get("auth-routes").(http.Handler)
	userRoutes := resources.Get("user-routes").(http.Handler)
	userRepoRoutes := resources.Get("user-repo-routes").(http.Handler)
	repositoryRoutes := resources.Get("repositories-routes").(http.Handler)
	captureRoutes := resources.Get("capture-routes").(http.Handler)
	branchRoutes := resources.Get("branch-routes").(http.Handler)
	multipostRoutes := resources.Get("multipost-routes").(http.Handler)

	r.Mount("/auth/", authRoutes)
	r.Mount("/users/", userRoutes)
	r.Route("/user/", func(r chi.Router) {
		r.Mount("/repos/", userRepoRoutes)
	})
	r.Mount("/repositories/", repositoryRoutes)
	r.Mount("/captures/", captureRoutes)
	r.Mount("/branches/", branchRoutes)
	r.Mount("/multipost/", multipostRoutes)

	return r
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	app := bastion.New(bastion.Addr(cfg.ADDR))
	app.APIRouter.Mount("/", router(cfg.Resources))
	app.RegisterOnShutdown(cfg.OnShutdown)
	if err := app.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
