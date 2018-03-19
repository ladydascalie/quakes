package handlers

import (
	"github.com/ladydascalie/earthquakes/app"
)

// RegisterRoutes sets up the routes against echo.
func RegisterRoutes(server *app.Server) {
	api := server.Echo.Group("api/v1")

	api.GET("/feed", FeedHandler(server))
	api.GET("/points", PointsHandler(server))

	// setup public group
	public := server.Echo.Group("")
	public.GET("/", ListHandler(server))
	public.GET("/:locale/:id", SingleHandler(server))
	public.GET("/display", DisplayHandler(server))
}
