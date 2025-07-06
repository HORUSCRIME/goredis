package utils

import (
	"log"
)

func HandlePanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}