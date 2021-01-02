package usershandler

import (
	"context"
	"time"

	"github.com/Thiti-Dev/AITTTY/database"
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/Thiti-Dev/AITTTY/models"
	encryptor "github.com/Thiti-Dev/AITTTY/packages/bcrypt"
	"github.com/Thiti-Dev/AITTTY/packages/jwt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StatusCheckAPI -> is for checking if the server is properly running at the moment
func StatusCheckAPI(c *fiber.Ctx) error {
	return helpers.ResponseMsg(c, 200, "API is up and running", nil)
}

// CreateAccount -> is the (fn) that will be creating the user account into the database
func CreateAccount(c *fiber.Ctx) error {
	users := new(models.Accounts) // using new will return an address of created struct
	var ctx = context.Background()

	db := database.GetDatabaseInstance()

	// Stamping the time
	users.CreatedAt = time.Now()
	users.UpdatedAt = time.Now()
	// ────────────────────────────────────────────────────────────────────────────────

	if err := c.BodyParser(users); err != nil {
		return helpers.ResponseMsg(c, 400, err.Error(), nil)
	} else {

		isValid, errorsData := helpers.ValidateStructAndGetErrorMsg(users)

		if !isValid{
			return helpers.ResponseMsg(c, 400, "Validation Errors", errorsData) 
		}


		//Check if this email already exist in db or not
		q := bson.M{"email": users.Email}
		result := db.Collection("users").FindOne(ctx, q)
		if result.Err() != nil {
			// To ensure if document isn't exist
			if result.Err() == mongo.ErrNoDocuments  { 

				// Encrypt the password
				hashPwd, _ := encryptor.HashPassword(users.Password)
				users.Password = hashPwd // replacing plain pwd with encrypted one
				// • • • • •


				//asynchronize task -> passing context to insertion task
				if r, err := db.Collection("users").InsertOne(ctx, users); err != nil {
					return helpers.ResponseMsg(c, 500, "Inserted data unsuccesfully", err.Error())
				} else {
					// Anonymous struct
					respData := struct {
						CreatedID    interface{}
						Username	string
						Email 		string
						Password    string
					}{
						r.InsertedID,
						users.Username,
						users.Email,
						hashPwd,
					}
					return helpers.ResponseMsg(c, 200, "Inserted data succesfully", respData)
				}
			}else{
				return helpers.ResponseMsg(c, 400, "Somethings went wrong", result.Err().Error())
			}
		}else{
			return helpers.ResponseMsg(c, 400, "This email is already existed", nil) 
		}
	}
}

// SignInWithCredential -> is the (fn) for signing-in the user and return the token for further usage
func SignInWithCredential(c *fiber.Ctx) error {
	//arbitData := make(map[string]interface{}) // -> ignore this (validation still need so i use struct instead)
	cred := new(models.AccountsSignIn)
	if err := c.BodyParser(cred); err != nil {
		return helpers.ResponseMsg(c, 400, err.Error(), nil)
	}

	// Validation
	isValid, errorsData := helpers.ValidateStructAndGetErrorMsg(cred)
	if !isValid{
		return helpers.ResponseMsg(c, 400, "Validation Errors", errorsData) 
	}
	// ────────────────────────────────────────────────────────────────────────────────

	// ─── CREDENTIAL MATCHING PHASE ──────────────────────────────────────────────────
	var ctx = context.Background()

	db := database.GetDatabaseInstance()

	q := bson.M{"email": cred.Email}
	userData := models.Accounts{}
	result := db.Collection("users").FindOne(ctx, q)
	result.Decode(&userData)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments  { 
			return helpers.ResponseMsg(c, 400, "email or password is incorrect", nil)
		}else{
			// any other error
			return helpers.ResponseMsg(c, 400, "Somethings went wrong", result.Err().Error())
		}
	}

	// checking if password are the same between plain&crypted
	if !encryptor.CheckPasswordHash(cred.Password,userData.Password) {
		return helpers.ResponseMsg(c, 400, "email or password is incorrect", nil)
	}
	 
	// ────────────────────────────────────────────────────────────────────────────────


	// ─── SIGN THE TOKEN FOR USER ────────────────────────────────────────────────────
	signedToken := jwt.GetSignedTokenFromData(models.RequiredDataToClaims{
		Username: userData.Username,
		Email: userData.Email,
		ID: userData.ID,
	})
	// ────────────────────────────────────────────────────────────────────────────────
	arbitResp := map[string]interface{}{
		"status": "success",
		"token": signedToken,
	}
	return c.Status(200).JSON(arbitResp)
	
}