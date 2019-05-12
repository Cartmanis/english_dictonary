package crypto

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"time"
)

var (
	salt = os.Getenv("SALT")
)

func GetHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func CreateToken() (string, error) {
	fmt.Println("salt: ", salt)
	h := md5.New()
	_, err := io.WriteString(h, salt+time.Now().String())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
