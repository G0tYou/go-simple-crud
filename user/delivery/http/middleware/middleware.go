package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// data-struct for middleware
type GoMiddleware struct {
}

// handle the CORS
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// initialize middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("B0mb45Tic"),
})
