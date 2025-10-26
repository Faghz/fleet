package request

type Login struct {
	Email        string `json:"email" validate:"required,email" example:"hi@example.com"`
	Password     string `json:"password" validate:"required" example:"IsThisAPassword?12345#@!"`
	IsRememberMe bool   `json:"isRememberMe" example:"true"`
}
