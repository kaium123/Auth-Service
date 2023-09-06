package utils

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func CopyStructToStruct(input interface{}, output interface{}) error {
	if byteData, err := json.Marshal(input); err == nil {
		if err := json.Unmarshal(byteData, &output); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func ComparePassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
