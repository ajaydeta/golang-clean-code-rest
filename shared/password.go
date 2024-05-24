package shared

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"os"
)

const PasswordSalt = "kerjakerjakerjapokoknyakerjasampaikerjabagus"

func EncryptPassword(password string) (string, error) {
	var (
		salted        = password + PasswordSalt
		passwordBytes = []byte(salted)
		cost          = bcrypt.DefaultCost
	)

	hash, errHash := bcrypt.GenerateFromPassword(passwordBytes, cost)
	if errHash != nil {
		return "", errHash
	}
	return string(hash), nil
}

func ComparePassword(storedHash, userPassword string) (bool, error) {
	var (
		salted = userPassword + PasswordSalt
	)

	godotenv.Load()
	masterPassword := os.Getenv("MASTER_PASSWORD") + PasswordSalt
	hashedMasterPassword, errhash := bcrypt.GenerateFromPassword([]byte(masterPassword), bcrypt.DefaultCost)
	if errhash != nil {
		return false, errhash
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(salted))
	if err != nil {
		err = bcrypt.CompareHashAndPassword(hashedMasterPassword, []byte(salted))
		if err != nil {
			return false, fmt.Errorf("invalid [password]")
		}
	}

	return true, nil
}
