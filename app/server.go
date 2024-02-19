package app

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ladydascalie/quakes/config"
	"golang.org/x/crypto/acme/autocert"
)

// Server defines the server and dependencies
// required for the app
type Server struct {
	Echo *echo.Echo
	DB   *mgo.Database
}

// NewServer returns a new instance of the server
func NewServer(db *mgo.Database, echo *echo.Echo) *Server {
	return &Server{
		Echo: echo,
		DB:   db,
	}
}

// Start the server with the provided config
func (s *Server) Start() {
	if config.Env == "prod" {
		// setup pre router middleware
		s.Echo.Pre(middleware.HTTPSRedirect())
		s.Echo.AutoTLSManager.HostPolicy = autocert.HostWhitelist("quakes.cable.fyi")
		s.Echo.AutoTLSManager.Cache = autocert.DirCache(".cache")
		s.Echo.StartAutoTLS(":443")
	} else {
		log.Fatalln(s.Echo.Start(config.AppPort))
	}
}
