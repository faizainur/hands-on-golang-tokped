package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
