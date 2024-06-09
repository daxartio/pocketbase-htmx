package app

import (
	"github.com/daxartio/pocketbase-htmx/internal/htmxutil"
	"github.com/daxartio/pocketbase-htmx/view"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

func Profile(c echo.Context) error {
	var record *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
	return htmxutil.Render(c, 200, view.Profile(record))
}
