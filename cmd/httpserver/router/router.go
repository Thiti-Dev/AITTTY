package router

import (
	"github.com/Thiti-Dev/AITTTY/handler/usershandler"
	"github.com/gofiber/fiber/v2"
)

// RoutesSetup -> A function that is used for setting up the endpoint route
func RoutesSetup(app *fiber.App){
	//TODO applied logger to the route
	api := app.Group("/api")

	api.Get("/",usershandler.StatusCheckAPI)

	//Users API related
	//userAPI := api.Group("/users")

}