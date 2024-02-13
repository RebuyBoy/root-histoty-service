package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Auth(keyword string, l *logrus.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(c echo.Context) error {
			tokenCookie, err := c.Cookie("access_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(keyword), nil
			})
			if err != nil {
				l.Info("failed to authorize:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			idStr, err := token.Claims.GetIssuer()
			if err != nil {
				l.Info("failed to authorize:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
			c.Set("userId", idStr)

			err = next(c)
			if err != nil {
				return err
			}
			return nil
		}

		return fn
	}
}
