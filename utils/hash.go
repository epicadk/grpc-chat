package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}
	return string(hashedpass), nil
}

func ComparePassword(hashedPass, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
