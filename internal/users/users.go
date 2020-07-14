package users

import (
	"context"
	"database/sql"

	"github.com/ohmpatel1997/twitter-graphql/graph/model"
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
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

func (user *User) Authenticate() (bool, error) {
	statement, err := database.Db.Prepare("select password from user WHERE user_name = ? OR email = ?")

	if err != nil {
		log.Println(err)
		return false, err
	}
	row := statement.QueryRow(user.Username, user.Email)

	var hashedPassword string
	err = row.Scan(&hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Println(err)
		return false, err
	}

	return CheckPasswordHash(user.Password, hashedPassword), nil
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (user *User) FetchAllTweetsOfUser(ctx context.Context) ([]*model.Tweet, error) {
	statement, err := database.Db.Prepare("select t.t_id, t.u_id, t.created_on, t.content from tweet t INNER JOIN user u ON t.u_id = u.u_id WHERE u.user_name = ? or u.u_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := statement.QueryContext(ctx, user.Username, user.ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	hasNext := rows.Next()
	var tweets []*model.Tweet
	for ; hasNext; hasNext = rows.Next() {
		tweet := model.Tweet{}
		if err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.CreatedOn, &tweet.Content); err != nil {
			log.Println(err)
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}
