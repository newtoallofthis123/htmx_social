package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GetEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	env := &Env{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		DB:       os.Getenv("DB"),
		Host:     os.Getenv("HOST"),
	}

	return env
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func MatchPasswords(toCheck string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(toCheck))
	return err == nil
}

func RanHash(limit int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var slug string

	for i := 0; i < limit; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			fmt.Println("Error generating random index:", err)
			return ""
		}
		slug += string(characters[randomIndex.Int64()])
	}

	return slug
}
