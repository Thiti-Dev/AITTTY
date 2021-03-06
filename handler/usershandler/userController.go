package usershandler

import (
	"context"
	"log"
	"time"

	"github.com/Thiti-Dev/AITTTY/database"
	"github.com/Thiti-Dev/AITTTY/helpers"
	"github.com/Thiti-Dev/AITTTY/models"
	encryptor "github.com/Thiti-Dev/AITTTY/packages/bcrypt"
	"github.com/Thiti-Dev/AITTTY/packages/jwt"
	jwtGOImplementedVer "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetAccounts -> getting all of the accounts existed in the database
func GetAccounts(c *fiber.Ctx) error {
	var ctx = context.Background()
	db := database.GetDatabaseInstance()

	// Also excluding the password field using projection
	csr, err := db.Collection("users").Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"password": 0,
	}))
	if err != nil {
		// Doesn't need to send back the error response -> fatal instead (if at somepoint it can't be execute this opreation -> os.exit is the neat way to handle)
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Accounts, 0)
	for csr.Next(ctx) {
		var row models.Accounts
		err := csr.Decode(&row)
		if err != nil {
			// Doesn't need to send back the error response -> fatal instead (if at somepoint it can't be execute this opreation -> os.exit is the neat way to handle)
			
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return helpers.ResponseMsg(c, 200, "Get Data Successfully", result)
}

// GetAccountByID -> get a single account information by ID
func GetAccountByID(c *fiber.Ctx) error {
	var ctx = context.Background()
	db := database.GetDatabaseInstance()
	_id := c.Params("id")

	if docID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return helpers.ResponseMsg(c, 400, "Get Data unsuccessfully (invalid ID type)", err.Error())
	} else {
		q := bson.M{"_id": docID}
		users := models.Accounts{}
		result := db.Collection("users").FindOne(ctx, q)
		result.Decode(&users)
		if result.Err() != nil {
			return helpers.ResponseMsg(c, 400, "Get Data unsuccesfully",result.Err().Error())
		} else {
			return helpers.ResponseMsg(c, 200, "Get Data Successfully", users)
		}
	}
}
// AuthorizedCheck -> A check if user is already authorized or not
func AuthorizedCheck(c *fiber.Ctx) error {
	// def name is user (if context key leaves as default in fiber/jwt)
	// Make use of jwtGOImplementedVer from  form3tech-oss because fiber/jwt supports the libs only from that source (satisfied types there)
	userData := c.Locals("user").(*jwtGOImplementedVer.Token) // cast the value -> because default is undefined so the next line would be blank (lack of knowledge) (need to be pointer cuz it's struct type)
	claims := userData.Claims.(jwtGOImplementedVer.MapClaims)
	username := claims["username"].(string)
	return c.SendString("Welcome " + username)
}