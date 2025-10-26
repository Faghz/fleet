package response

import (
	"net/http"
)

const (
	userNotFoundMessage        = "User record not found"
	userEmailAlreadyUsed       = "Email Already Used"
	userInvalidUserIDError     = "Invalid user id"
	InternalServerErrorMessage = "Internal Server Error"
)

var (
	ErrorUserServiceInternalServerError   = GenerateFailure(http.StatusInternalServerError, InternalServerErrorMessage, InternalServerErrorMessage)
	ErrorUserDatabaseInternalError        = GenerateFailure(http.StatusInternalServerError, InternalServerErrorMessage, InternalServerErrorMessage)
	ErrorUserDatabaseUserNotFound         = GenerateFailure(http.StatusNotFound, userNotFoundMessage, userNotFoundMessage)
	ErrorUserDatabaseUserEmailAlreadyUsed = GenerateFailure(http.StatusConflict, userEmailAlreadyUsed, userEmailAlreadyUsed)
	ErrorInvalidUserID                    = GenerateFailure(http.StatusBadRequest, userInvalidUserIDError, userInvalidUserIDError)
	ErrorInvalidEmailOrPassword           = GenerateFailure(http.StatusBadRequest, "Invalid email or password", "Invalid email or password")
	ErrorUserDatabaseUserNotActive        = GenerateFailure(http.StatusBadRequest, "IPlease contact the administrator for more details", "Please contact the administrator for more details")
)

type UserDetail struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	OrgID     string `json:"org_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
