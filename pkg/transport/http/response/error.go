package response

import (
	"net/http"
)

var (
	ErrorUserBadRequest      = GenerateFailure(http.StatusBadRequest, "Bad Request", "Bad Request")
	ErrorInternalServerError = GenerateFailure(http.StatusInternalServerError, "Internal Server Error", "Internal Server Error")
	InvalidToken             = GenerateFailure(http.StatusUnauthorized, "Session Not Valid", "Session Not Valid")
	ErrorUnAuthorized        = GenerateFailure(http.StatusUnauthorized, "Unauthorized", "Unauthorized")
	ErrorInvalidArgument     = GenerateFailure(http.StatusBadRequest, "Invalid Argument", "Invalid Argument")
	ErrorNotFound            = GenerateFailure(http.StatusNotFound, "Not Found", "Not Found")
)
