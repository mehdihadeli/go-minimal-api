package minimalapi

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"net/http"
)

type WebApplicationBuilder struct {
	port         int
	dependencies []interface{}
	options      []fx.Option
	lifecycle    fx.Lifecycle
	httpHandlers []echo.HandlerFunc
	logger       slog.Logger
	config       *viper.Viper
}

const defaultPort = 8080

func NewWebApplicationBuilder() *WebApplicationBuilder {
	return &WebApplicationBuilder{}
}

func (builder *WebApplicationBuilder) AddDependency(deps ...interface{}) {
	builder.dependencies = append(builder.dependencies, deps...)
}

func (builder *WebApplicationBuilder) WithLogger(logger slog.Logger) *WebApplicationBuilder {
	builder.logger = logger
	return builder
}

// WithConfig sets up the configuration from a YAML file.
func (b *WebApplicationBuilder) WithConfig(config *viper.Viper) *WebApplicationBuilder {
	b.config = config
	return b
}
func (builder *WebApplicationBuilder) WithOptions(options ...fx.Option) *WebApplicationBuilder {
	builder.options = append(builder.options, options...)
	return builder
}

func (builder *WebApplicationBuilder) WithLifecycle(lifecycle fx.Lifecycle) *WebApplicationBuilder {
	builder.lifecycle = lifecycle
	return builder
}

func (builder *WebApplicationBuilder) WithPort(port int) *WebApplicationBuilder {
	builder.port = port
	return builder
}

func (builder *WebApplicationBuilder) AddHTTPHandler(handler echo.HandlerFunc) *WebApplicationBuilder {
	builder.httpHandlers = append(builder.httpHandlers, handler)

	return builder
}

func (builder *WebApplicationBuilder) Build() (*WebApplication, error) {
	// convert dependencies to fx.Option
	opts := make([]fx.Option, len(builder.dependencies))
	for i, dep := range builder.dependencies {
		opts[i] = fx.Provide(dep)
	}

	// Create the echo server
	e := echo.New()

	app := fx.New(
		fx.Provide(
			builder.provideServer,
			func() *echo.Echo {
				return e
			},
		),
		fx.Invoke(builder.invokeEchoHandlers),
		fx.Provide(builder.provideServer),
		fx.Options(opts...), // pass as variadic parameter
	)

	webApp := &WebApplication{
		app: app,
	}

	webApp.logger = builder.logger

	if builder.options != nil {
		webApp.options = builder.options
	}

	if len(builder.httpHandlers) > 0 {
		webApp.httpHandlers = builder.httpHandlers
	}

	if builder.port == 0 {
		builder.port = defaultPort
	}

	return &WebApplication{
		app:    app,
		server: e,
	}, nil
}

func (builder *WebApplicationBuilder) provideServer(lifecycle fx.Lifecycle, e *echo.Echo) *http.Server {
	server := e.Server
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil {
					e.Logger.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}
			return nil
		},
	})
	return server
}

func (builder *WebApplicationBuilder) invokeEchoHandlers(e *echo.Echo) {
	// Add your echo route handlers here
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
