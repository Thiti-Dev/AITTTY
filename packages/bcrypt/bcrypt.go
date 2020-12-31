package packages

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword -> is the (fn) that will hash your given password and return if it successes or not
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}


// CheckPasswordHash -> is the (fn) that will check if the plain-txt-password & crypted-password is matched or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}