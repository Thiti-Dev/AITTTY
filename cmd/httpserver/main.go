package main

import (
	"log"

	"github.com/Thiti-Dev/AITTTY/cmd/httpserver/router"
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	//Pkg initializer
	helpers.InitializeTranslator()
	// ────────────────────────────────────────────────────────────────────────────────


	app := fiber.New()
	app.Use(logger.New()) // appiled logger middleware
	router.RoutesSetup(app)

	log.Fatal(app.Listen(":3000"))
}