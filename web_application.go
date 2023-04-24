package minimalapi

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type WebApplication struct {
	app          *fx.App
	server       *echo.Echo
	logger       slog.Logger
	options      []fx.Option
	httpHandlers []echo.HandlerFunc
}

func (app *WebApplication) RegisterHTTPHandler(handlerFuncs ...echo.HandlerFunc) {
	app.httpHandlers = append(app.httpHandlers, handlerFuncs...)
}

func (app *WebApplication) Run() {
	app.logger.Info("Starting server...")
	app.server.Logger.Fatal(app.server.Start(":8080"))
}

func (app *WebApplication) loggerMiddleware(next echo.HandlerFunc) *WebApplication {
	return app
}

func (app *WebApplication) GET(path string, handler echo.HandlerFunc) *WebApplication {
	app.server.GET(path, handler)
	return app
}

func (app *WebApplication) POST(path string, handler echo.HandlerFunc) *WebApplication {
	app.server.POST(path, handler)
	return app
}
