package auth

import (
	"net/http"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Login func(error) templ.Component
type LoginForm func(error) templ.Component

type LoginFormValue struct {
	username string
	password string
}

func (lfv LoginFormValue) Validate() error {
	return validation.ValidateStruct(&lfv,
		validation.Field(&lfv.username, validation.Required, validation.Length(3, 50)),
		validation.Field(&lfv.password, validation.Required),
	)
}

func getLoginFormValue(c echo.Context) LoginFormValue {
	return LoginFormValue{
		username: c.FormValue("username"),
		password: c.FormValue("password"),
	}
}

func (a *Auth) RegisterLoginRoutes(e *core.ServeEvent, group echo.Group) {
	group.GET(a.cfg.LoginPath, func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, a.cfg.RedirectPath)
		}

		return Render(c, 200, a.cfg.Login(nil))
	})

	group.POST(a.cfg.LoginPath, func(c echo.Context) error {
		form := getLoginFormValue(c)
		err := form.Validate()

		if err == nil {
			err = a.Login(e, c, form.username, form.password)
		}

		if err != nil {
			component := HtmxRender(
				c,
				func() templ.Component { return a.cfg.LoginForm(err) },
				func() templ.Component { return a.cfg.Login(err) },
			)
			return Render(c, 200, component)
		}

		return HtmxRedirect(c, a.cfg.RedirectPath)
	})

	group.POST(a.cfg.LogoutPath, func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     a.cfg.AuthCookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   -1,
		})

		return HtmxRedirect(c, a.cfg.RedirectLoginPath)
	})
}
