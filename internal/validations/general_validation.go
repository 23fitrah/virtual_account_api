package validations

type UserValidation struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Channel  string `json:"channel" validate:"required"`
}
