package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ladydascalie/earthquakes/app"
	"github.com/ladydascalie/earthquakes/app/handlers"
	"github.com/ladydascalie/earthquakes/bot"
	"github.com/ladydascalie/earthquakes/config"
	"github.com/ladydascalie/earthquakes/db/mongo"
	"github.com/ladydascalie/earthquakes/templates"
	"github.com/ladydascalie/earthquakes/workers"
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
	go bot.Begin()

	// Create the server and start it
	server := app.NewServer(db, e)

	// register the routes
	handlers.RegisterRoutes(server)

	// start the server
	server.Start()
}
