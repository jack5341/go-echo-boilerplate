package middlewares

import (
	"backend/internal/svc"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
)

func Trace(s *svc.ServiceContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := *s.Tracer
			_, span := tracer.Start(c.Request().Context(), "middleware.Trace")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", c.Request().Method))
			span.SetAttributes(attribute.String("http.route", c.Request().URL.Path))
			span.SetAttributes(attribute.String("http.request.id", c.Response().Header().Get(echo.HeaderXRequestID)))
			span.SetAttributes(attribute.String("client.ip", c.RealIP()))
			span.SetAttributes(attribute.String("http.user_agent", c.Request().UserAgent()))
			span.SetAttributes(attribute.Int64("http.request.body.size", c.Request().ContentLength))
			return next(c)
		}
	}
}
