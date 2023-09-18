package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"backend/internal/handler"
	"backend/internal/svc"
	"backend/pkg/config"
	"backend/pkg/database"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"

	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Add a commentx to the main function
// @title Go-Boilerplate
// @version 1.0
// @description This is Go-Boilerplate Backend Api Documentation
//
//	<a href="https://swagger.io">here</a>.
//
// @host go-boilerplate.nedim-akar.cloud
// @BasePath /
func main() {
	// enable multi-span attributes
	bsp := honeycomb.NewBaggageSpanProcessor()

	// use honeycomb distro to setup OpenTelemetry SDK
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithSpanProcessor(bsp),
	)
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()

	var tracer = otel.GetTracerProvider().Tracer("go-boilerplate")

	e := echo.New()
	cfg := config.InitConfig()

	conn, _ := database.ConnectDB()

	err = conn.AutoMigrate()

	if err != nil {
		e.Logger.Fatal(err)
	}

	env, err := strconv.ParseBool(cfg.APP.DEV)
	if err != nil {
		cfg.DevMode = true
	}

	cfg.DevMode = env

	if !cfg.DevMode {
		e.Use(middleware.BodyLimit("30MB"))
		e.Use(middleware.Decompress())
		e.Use(middleware.Gzip())
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())
		e.Use(otelecho.Middleware("go-boilerplate"))
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(50)))
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"go-boilerplate.nedim-akar.cloud"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
		e.Use(otelecho.Middleware("go-boilerplate"))
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.GET("/buildz", func(c echo.Context) error {
			ctx := c.Request().Context()
			_, span := tracer.Start(ctx, "handler.buildz")
			span.SetName("buildz")
			info, _ := debug.ReadBuildInfo()
			return c.JSON(http.StatusOK, info)
		})

		e.GET("/debug", func(c echo.Context) error {
			info, _ := debug.ReadBuildInfo()
			ctx := c.Request().Context()
			_, span := tracer.Start(ctx, "handler.debug")
			span.SetName("debug")
			return c.JSON(http.StatusOK, info)
		})
	}

	// Health Check
	e.GET("/healthz", func(c echo.Context) error {
		ctx := c.Request().Context()
		_, span := tracer.Start(ctx, "handler.healthz")
		span.SetName("healthz")
		return c.JSON(http.StatusOK, echo.Map{
			"message": "everything is ok",
		})
	})

	serviceCtx := svc.NewServiceContext(cfg, database.DB, e, &tracer)
	handler.RegisterHandlers(serviceCtx)

	s := http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              ":" + cfg.APP.PORT,
		Handler:           hnynethttp.WrapHandler(e),
	}

	fmt.Println("server does work on " + cfg.APP.PORT)
	e.Logger.Fatal(s.ListenAndServe())
}
