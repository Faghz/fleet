package request

type RegisterUserRequest struct {
	Name            string `json:"name" validate:"required" example:"John"`
	Email           string `json:"email" validate:"required,email" example:"hi@example.com"`
	Password        string `json:"password" validate:"required" example:"IsThisAPassword?12345#@!"`
	ConfirmPassword string `json:"confirmPassword" validate:"required" example:"IsThisAPassword?12345#@!"`
}
