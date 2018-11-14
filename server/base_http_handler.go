package main

import "net/http"

// Contains common helper methods.
type baseHandler struct{}

func (handler *baseHandler) sendError(response http.ResponseWriter, statusCode int) {
	response.WriteHeader(statusCode)
	response.Write([]byte(http.StatusText(statusCode)))
}
