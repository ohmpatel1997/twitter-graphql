package users

import (
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
	"log"
)

type User struct {
	ID        int    `json:"u_id"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`
	FirstName string `json:"f_name"`
	LastName  string `json:"l_name"`
	Email     string `json:"email"`
	Deleted   bool   `json:"deleted"`
}

func (user *User) Save() (int64, error) {
	statement, err := database.Db.Prepare("INSERT INTO user(f_name,l_name,email,password,user_name) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		return -1, err
	}

	res, err := statement.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Username)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	log.Print("User inserted!")
	return id, nil
}
