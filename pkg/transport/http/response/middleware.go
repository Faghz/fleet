package response

import (
	"net/http"
)

var (
	ErrorMiddlewareInternalError    = GenerateFailure(http.StatusInternalServerError, "Internal Server Error", "Internal Server Error")
	ErrorMiddlewareMissingToken     = GenerateFailure(http.StatusUnauthorized, "No authorization token provided", "No authorization token provided")
	ErrorMiddlewareTokenExpired     = GenerateFailure(http.StatusUnauthorized, "Token expired", "Token expired")
	ErrorMiddlewareNotAuthorized    = GenerateFailure(http.StatusUnauthorized, "Not authorized", "Not authorized")
	ErrorMiddlewareNotAuthenticated = GenerateFailure(http.StatusUnauthorized, "Not authenticated", "Not authenticated")
	ErrorMiddlewareTokenInvalid     = GenerateFailure(http.StatusUnauthorized, "Invalid token", "Invalid token")
)
