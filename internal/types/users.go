package types

import "time"

type UserStore interface {
	CreateUser(data RegisterPayload) error
	CreateResetToken(userId int) (string, error)
	GetUserCredentials(name string) (Credentials, error)
	GetUserFromEmail(email string) (int, error)
	GetDataFromResetToken(token string) (ResetTokenData, error)
	UpdateUserPassword(userId int, password string) error
	DeleteUser(userId int) error
}

type Login struct {
	Token string `json:"token"`
}

type RegisterPayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ForgotPasswordPayload struct {
	Email string `json:"email"`
}

type ResetPasswordPayload struct {
	Password string `json:"password"`
}

type Credentials struct {
	Id       int
	Password string
}

type ResetTokenData struct {
	UserId     int
	Expiration time.Time
}
