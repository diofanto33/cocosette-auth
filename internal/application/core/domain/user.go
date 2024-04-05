package domain

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(email, password string) User {
	return User{
		Email:    email,
		Password: password,
	}
}
