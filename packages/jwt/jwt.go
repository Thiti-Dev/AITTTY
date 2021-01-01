package jwt

import (
	"fmt"

	"github.com/Thiti-Dev/AITTTY/models"
	"github.com/dgrijalva/jwt-go"
)

// GetSignedTokenFromData -> is a (fn) that will sign a token with the given data
func GetSignedTokenFromData(data models.RequiredDataToClaims) string{
	claims := models.CustomClaims{
		Username: data.Username,
		Email: data.Email,
		ID: data.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "aittty.io",
		},
	}
	// Create an unsigned token from the claims above
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token ->  preferably at least 256 bits in length (in-production xD)
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	if err != nil {
		fmt.Println(err)
	}
	return signedToken
}