package svc

import (
	"backend/pkg/config"

	"go.opentelemetry.io/otel/trace"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Configuration
	DB     *gorm.DB
	Echo   *echo.Echo
	Tracer *trace.Tracer
}

func NewServiceContext(c config.Configuration, d *gorm.DB, e *echo.Echo, t *trace.Tracer) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     d,
		Echo:   e,
		Tracer: t,
	}
}
