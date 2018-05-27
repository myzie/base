package env

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pavel-kiselyov/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"
)

// SetupEcho creates a new Echo instance for HTTP
func SetupEcho() (*echo.Echo, error) {

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(logrusmiddleware.Hook())
	e.Use(middleware.Recover())

	return e, nil
}
