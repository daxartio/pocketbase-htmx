package main

import (
	"log"

	"github.com/daxartio/pocketbase-htmx/internal/app"
	"github.com/daxartio/pocketbase-htmx/internal/auth"
	"github.com/daxartio/pocketbase-htmx/view"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	pb := pocketbase.New()

	auth := auth.New(func(cfg *auth.Config) {
		cfg.TableName = "users"
		cfg.AuthCookieName = "Auth"
		cfg.LoginPath = "/login"
		cfg.RegisterPath = "/register"
		cfg.LogoutPath = "/logout"
		cfg.RedirectPath = "/"
		cfg.RedirectLoginPath = "/auth/login"
		cfg.Login = view.Login
		cfg.LoginForm = view.LoginForm
		cfg.Register = view.Register
		cfg.RegisterForm = view.RegisterForm
	})

	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Static("/public", "public")

		authGroup := e.Router.Group("/auth", auth.LoadAuthContextFromCookieMiddleware(pb))
		auth.RegisterLoginRoutes(e, *authGroup)
		auth.RegisterRegisterRoutes(e, *authGroup)

		appGroup := e.Router.Group("", auth.LoadAuthContextFromCookieMiddleware(pb), auth.AuthGuardMiddleware())
		appGroup.GET("", app.ProfileGet)

		return nil
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
