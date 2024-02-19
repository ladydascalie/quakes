package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ladydascalie/quakes/app"
	"github.com/ladydascalie/quakes/app/handlers"
	"github.com/ladydascalie/quakes/config"
	"github.com/ladydascalie/quakes/db/mongo"
	"github.com/ladydascalie/quakes/templates"
	"github.com/ladydascalie/quakes/workers"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	// setup the app config
	config.Setup()

	// create a new instance of Echo and configure it
	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.New()

	// configure middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	// configure static assets folder
	e.Static("static", "static")

	// Create a new database connection
	db := mongo.Begin()

	go workers.Run(db)
	// go bot.Begin()

	// Create the server and start it
	server := app.NewServer(db, e)

	// register the routes
	handlers.RegisterRoutes(server)

	// start the server
	server.Start()
}
