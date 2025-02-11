package main

import "log"

func main() {
	s, err := NewServer()
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	s.Start()
}
