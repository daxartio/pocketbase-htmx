package auth

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func (a *Auth) LoadAuthContextFromCookieMiddleware(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenCookie, err := c.Request().Cookie(a.cfg.AuthCookieName)
			if err != nil {
				return next(c)
			}

			token := tokenCookie.Value
			record, err := app.Dao().FindAuthRecordByToken(
				token,
				app.Settings().RecordAuthToken.Secret,
			)

			if err != nil {
				return next(c)
			}

			c.Set(apis.ContextAuthRecordKey, record)
			return next(c)
		}
	}
}

func (a *Auth) AuthGuardMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record := c.Get(apis.ContextAuthRecordKey)

			if record == nil {
				return c.Redirect(302, a.cfg.RedirectLoginPath)
			}

			return next(c)
		}
	}
}
