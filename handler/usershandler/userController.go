package usershandler

import (
	"context"
	"log"
	"time"

	"github.com/Thiti-Dev/AITTTY/database"
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/Thiti-Dev/AITTTY/models"
	"github.com/gofiber/fiber/v2"
)

// StatusCheckAPI -> is for checking if the server is properly running at the moment
func StatusCheckAPI(c *fiber.Ctx) error {
	return helpers.ResponseMsg(c, 200, "API is up and running", nil)
}

// CreateAccount -> is the (fn) that will be creating the user account into the database
func CreateAccount(c *fiber.Ctx) error {
	users := new(models.Accounts)
	var ctx = context.Background()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Stamping the time
	users.CreatedAt = time.Now()
	users.UpdatedAt = time.Now()
	// ────────────────────────────────────────────────────────────────────────────────

	if err := c.BodyParser(users); err != nil {
		return helpers.ResponseMsg(c, 400, err.Error(), nil)
	} else {
		//asynchronize task -> passing context to insertion task
		if r, err := db.Collection("users").InsertOne(ctx, users); err != nil {
			return helpers.ResponseMsg(c, 500, "Inserted data unsuccesfully", err.Error())
		} else {
			// Anonymous struct
			respData := struct {
				CreatedID    interface{}
				Username	string
				Email 		string
			}{
				r.InsertedID,
				users.Username,
				users.Email,
			}
			return helpers.ResponseMsg(c, 200, "Inserted data succesfully", respData)
		}
	}
}