package auth

import (
	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Register func(error) templ.Component
type RegisterForm func(error) templ.Component

type RegisterFormValue struct {
	username       string
	password       string
	passwordRepeat string
}

func (lfv RegisterFormValue) Validate() error {
	return validation.ValidateStruct(&lfv,
		validation.Field(&lfv.username, validation.Required, validation.Length(3, 50)),
		validation.Field(&lfv.password, validation.Required),
	)
}

func getRegisterFormValue(c echo.Context) RegisterFormValue {
	return RegisterFormValue{
		username:       c.FormValue("username"),
		password:       c.FormValue("password"),
		passwordRepeat: c.FormValue("passwordRepeat"),
	}
}

func (a *Auth) RegisterRegisterRoutes(e *core.ServeEvent, group echo.Group) {
	group.GET(a.cfg.RegisterPath, func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, a.cfg.RedirectPath)
		}

		return Render(c, 200, a.cfg.Register(nil))
	})

	group.POST(a.cfg.RegisterPath, func(c echo.Context) error {
		form := getRegisterFormValue(c)
		err := form.Validate()

		if err == nil {
			err = a.Register(e, c, form.username, form.password, form.passwordRepeat)
		}

		if err != nil {
			component := HtmxRender(
				c,
				func() templ.Component { return a.cfg.RegisterForm(err) },
				func() templ.Component { return a.cfg.Register(err) },
			)
			return Render(c, 200, component)
		}

		return HtmxRedirect(c, a.cfg.RedirectPath)
	})
}
