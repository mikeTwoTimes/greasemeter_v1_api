package types

type UserStore interface {
	CreateUser(data RegisterPayload) error
	GetUserCredentials(name string) (Credentials, error)
	GetUserByEmail(email string) (int, error)
	UserExists(userId int) (bool, error)
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
