package usershandler

import (
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/gofiber/fiber/v2"
)

// StatusCheckAPI -> is for checking if the server is properly running at the moment
func StatusCheckAPI(c *fiber.Ctx) error {
	return helpers.ResponseMsg(c, 200, "API is up and running", nil)
}