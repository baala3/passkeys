package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	s, err := NewServer()
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	s.Start()
}
