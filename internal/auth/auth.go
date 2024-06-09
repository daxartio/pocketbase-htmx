package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

type Option func(cfg *Config)

type Config struct {
	TableName         string
	AuthCookieName    string
	LoginPath         string
	RegisterPath      string
	LogoutPath        string
	RedirectPath      string
	RedirectLoginPath string
	Login             Login
	LoginForm         LoginForm
	Register          Register
	RegisterForm      RegisterForm
}

type Auth struct {
	cfg Config
}

func New(options ...Option) *Auth {
	cfg := &Config{
		TableName:         "users",
		AuthCookieName:    "Auth",
		LoginPath:         "/login",
		RegisterPath:      "/register",
		LogoutPath:        "/register",
		RedirectPath:      "",
		RedirectLoginPath: "",
		Login:             nil,
		LoginForm:         nil,
		Register:          nil,
		RegisterForm:      nil,
	}

	for _, fn := range options {
		fn(cfg)
	}

	return &Auth{cfg: *cfg}
}

type Users struct {
	models.Record
}

func (a *Auth) Login(e *core.ServeEvent, c echo.Context, username string, password string) error {
	user, err := e.App.Dao().FindAuthRecordByUsername(a.cfg.TableName, username)
	if err != nil {
		return fmt.Errorf("Login failed")
	}

	valid := user.ValidatePassword(password)
	if !valid {
		return fmt.Errorf("Login failed")
	}

	return a.setAuthToken(e.App, c, user)
}

func (a *Auth) Register(e *core.ServeEvent, c echo.Context, username string, password string, passwordRepeat string) error {
	user, _ := e.App.Dao().FindAuthRecordByUsername(a.cfg.TableName, username)
	if user != nil {
		return fmt.Errorf("username already taken")
	}

	if password != passwordRepeat {
		return fmt.Errorf("passwords don't match")
	}

	collection, err := e.App.Dao().FindCollectionByNameOrId(a.cfg.TableName)
	if err != nil {
		return err
	}

	newUser := models.NewRecord(collection)
	newUser.SetPassword(password)
	newUser.SetUsername(username)

	if err = e.App.Dao().SaveRecord(newUser); err != nil {
		return err
	}

	return a.setAuthToken(e.App, c, newUser)
}

func (a *Auth) setAuthToken(app core.App, c echo.Context, user *models.Record) error {
	s, tokenErr := tokens.NewRecordAuthToken(app, user)
	if tokenErr != nil {
		return fmt.Errorf("Login failed")
	}

	c.SetCookie(&http.Cookie{
		Name:     a.cfg.AuthCookieName,
		Value:    s,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}
