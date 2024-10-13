package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port       string
	gqlHandler handler.Server
	server     *echo.Echo
}

func NewServer(port int, gqlHandler handler.Server) *Server {
	return &Server{
		port:       fmt.Sprintf(":%d", port),
		gqlHandler: gqlHandler,
		server:     echo.New(),
	}
}

func (s *Server) Start(hasPlayground bool) error {
	s.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			// ignore health check, metrics
			// ignore HttpStatusCode: 200
			return strings.Contains(c.Path(), "health") ||
				strings.Contains(c.Path(), "metrics") ||
				c.Response().Status == http.StatusOK
		},
	}))
	s.server.Use(middleware.Recover())
	// GraphQL
	s.server.POST("/graphql", func(c echo.Context) error {
		s.gqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	// metrics
	mwConf := echoprometheus.MiddlewareConfig{
		Subsystem: "go-gql-sqlc-template",
		Skipper: func(c echo.Context) bool {
			// ignore health check, metrics
			return strings.Contains(c.Path(), "health") || strings.Contains(c.Path(), "metrics")
		},
	}
	s.server.Use(echoprometheus.NewMiddlewareWithConfig(mwConf))
	s.server.GET("/metrics", echoprometheus.NewHandler())

	if hasPlayground {
		playgroundHandler := playground.Handler("GraphQL playground", "/graphql")
		s.server.GET("/", func(c echo.Context) error {
			playgroundHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		})
	}

	return s.server.Start(s.port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Server is used in test
func (s *Server) Server() *echo.Echo {
	return s.server
}
