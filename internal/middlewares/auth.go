package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"fmt"

	"backend/pkg/config"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat/go-jwx/jwk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// JWK represents the JSON Web Key structure
type JWK struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

// GetCognitoPublicKeys retrieves the public keys from Cognito using the region and user pool ID
func GetCognitoPublicKeys(ctx context.Context, region, userPoolID string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		_, span := tracer.Start(ctx, "helper.GetCognitoPublicKeys")
		jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
		set, err := jwk.FetchHTTP(jwksURL)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have string kid")
		}

		if key := set.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}

		return nil, errors.New("unable to find key")
	}
}

var tracer = otel.GetTracerProvider().Tracer("middleware.AuthValidator")

func AuthValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		ctx, span := tracer.Start(c.Request().Context(), "middelware.AuthValidator")
		defer span.End()

		if tokenString == "" {
			span.RecordError(errors.New("token couldn't be parsed"))
			span.SetAttributes(attribute.Key("AuthHeader").String(authHeader))
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		cfg := config.InitConfig()
		keyFunc := GetCognitoPublicKeys(ctx, cfg.AWS.REGIONS, cfg.AWS.COGNITO.USERPOOL_ID)
		token, err := jwt.Parse(tokenString, keyFunc)

		if err != nil {
			span.RecordError(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid access token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Failed to get token claims",
			})
		}

		userId, ok := claims["cognito:username"].(string)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Failed to get subject claim",
			})
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Access token has expired",
			})
		}

		span.SetAttributes(attribute.Key("user.id").String(userId))
		c.Request().Header.Add("user.id", userId)
		return next(c)
	}
}
