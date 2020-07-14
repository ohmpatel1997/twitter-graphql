package users

type User struct {
	ID        int    `json:"u_id"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`
	FirstName string `json:"f_name"`
	LastName  string `json:"l_name"`
	Email     string `json:"email"`
	Deleted   bool   `json:"deleted"`
}
