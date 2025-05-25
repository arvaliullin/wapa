package utils

import (
	"log"
)

func PanicOnError(err error, message string, args any) {
	if err != nil {
		log.Panicf(message, args)
	}
}

func LogOnError(err error, message string, args any) {
	if err != nil {
		log.Printf(message, args)
	}
}

func FatalOnError(err error, message string, args any) {
	if err != nil {
		log.Fatalf(message, args)
	}
}
