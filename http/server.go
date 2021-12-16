package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/music-gang/music-gang-api/app/service"
	"golang.org/x/crypto/acme/autocert"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

// ServerAPI is the main server for the API
type ServerAPI struct {
	ln net.Listener
	// server is the main server for the API
	server *http.Server

	// handler is the main handler for the API
	handler *echo.Echo

	// Addr Bind address for the server.
	Addr string
	// Domain name to use for the server.
	// If specified, server is run on TLS using acme/autocert.
	Domain string

	// JWTSecret is the secret used to sign JWT tokens.
	JWTSecret string

	// Services used by HTTP handler.
	AuthService service.AuthService
	UserService service.UserService
	JWTService  service.JWTService
}

// NewServerAPI creates a new API server.
func NewServerAPI() *ServerAPI {

	jwtSecret := "secret"

	s := &ServerAPI{
		server:    &http.Server{},
		handler:   echo.New(),
		JWTSecret: jwtSecret,
	}

	// Set echo as the default HTTP handler.
	s.server.Handler = s.handler

	// Base Middleware
	s.handler.Use(middleware.Recover())

	// Register routes for the API v1.
	v1Group := s.handler.Group("/v1")
	s.registerRoutes(v1Group)

	return s
}

func (s *ServerAPI) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *ServerAPI) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Open validates the server options and start it on the bind address.
func (s *ServerAPI) Open() (err error) {

	if s.Domain != "" {
		s.ln = autocert.NewListener(s.Domain)
	} else {
		if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
			return err
		}
	}

	go s.server.Serve(s.ln)

	return nil
}

// Scheme returns the scheme used by the server.
func (s *ServerAPI) Scheme() string {
	if s.Domain != "" {
		return "https"
	}
	return "http"
}

// URL returns the URL for the server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *ServerAPI) URL() string {

	scheme, port := s.Scheme(), s.Port()

	domain := "localhost"

	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", scheme, domain)
	}

	return fmt.Sprintf("%s://%s:%d", scheme, domain, port)
}

// UseTLS returns true if the server is using TLS.
func (s *ServerAPI) UseTLS() bool {
	return s.Domain != ""
}

// registerRoutes registers all routes for the API.
func (s *ServerAPI) registerRoutes(g *echo.Group) {

	authGroup := g.Group("/auth")
	s.registerAuthRoutes(authGroup)
}

// registerAuthRoutes registers all routes for the API group auth.
func (s *ServerAPI) registerAuthRoutes(g *echo.Group) {
	g.POST("/login", s.AuthLogin)
	g.POST("/register", s.AuthRegister)
	g.GET("/refresh", nil)
	g.DELETE("/logout", s.AuthLogout)

	// oauth2 routes
	g.GET("/oauth2/:source/callback", nil)
}

// SuccessResponseJSON returns a JSON response with the given status code and data.
func SuccessResponseJSON(c echo.Context, httpCode int, data interface{}) error {
	return c.JSON(httpCode, data)
}
