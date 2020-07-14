package users

type InvalidUsernameOrPasswordError struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

func (m *InvalidUsernameOrPasswordError) Error() string {
	return "Authentication Failed..! Invalid username/email or password"
}
