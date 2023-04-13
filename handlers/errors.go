package handlers

import (
	"log"
)

func HandleErrors(options *HandlerOptions) {
	log.Println(options.Error)
	options.Status.Errors++
}
