package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Thiti-Dev/AITTTY/cmd/httpserver/router"
	"github.com/Thiti-Dev/AITTTY/database"
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	//Ensuring the repository
	if err := database.ConnectAndHoldTheConnection(); err != nil{
		log.Fatal(err)
	}
	// ────────────────────────────────────────────────────────────────────────────────


	//Pkg initializer
	helpers.InitializeTranslator()
	// ────────────────────────────────────────────────────────────────────────────────


	app := fiber.New()
	app.Use(logger.New()) // appiled logger middleware
	router.RoutesSetup(app)

	log.Fatal(app.Listen(":3000"))


	// ─── DISCONNECT A CONNTECTIN TO DATABASE ────────────────────────────────────────
	defer func() {
		err := database.GetDatabaseInstance().Client().Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB is closed.")
	}()
	// ────────────────────────────────────────────────────────────────────────────────

}