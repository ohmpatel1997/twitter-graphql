package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ohmpatel1997/twitter-graphql/common"
	"github.com/ohmpatel1997/twitter-graphql/graph/model"
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
	defer statement.Close()
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

	defer statement.Close()
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

	return common.CheckPasswordHash(user.Password, hashedPassword), nil
}

func (user *User) FetchAllTweetsOfUser(ctx context.Context) ([]*model.Tweet, error) {
	statement, err := database.Db.Prepare("select t.t_id, t.u_id, t.created_on, t.content from tweet t INNER JOIN user u ON t.u_id = u.u_id WHERE u.user_name = ? or u.u_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer statement.Close()
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

	if len(tweets) == 0 {
		log.Println("Could not able to find the tweets")
		return nil, fmt.Errorf("No tweets available for the user")
	}

	return tweets, nil
}

func (user *User) AddFollower(ctx context.Context, followerID int) (bool, error) {
	statement, err := database.Db.Prepare("select user_id, follower_id, active from follower where user_id = ? and follower_id = ?")

	if err != nil {
		log.Println(err)
		return false, err
	}

	defer statement.Close()

	res := statement.QueryRow(user.ID, followerID)

	var follower model.Relationship

	err = res.Scan(&follower.UserID, &follower.FollowerID, &follower.Active)

	if err != nil && err != sql.ErrNoRows { // proceed if no rows found
		log.Println(err)
		return false, err
	}

	if follower.Active != nil && *follower.Active { // return if user already followes
		log.Println("Relationship Already Exist")
		err := fmt.Errorf("User Already follows")
		return false, err
	}

	if follower.Active != nil && !*follower.Active { // update the active flag
		statement, err = database.Db.Prepare("update follower set active = true where user_id = ? and follower_id = ?")
		if err != nil {
			log.Println(err)
			return false, err
		}

		defer statement.Close()
		_, err := statement.Exec(user.ID, followerID)

		if err != nil {
			log.Println(err)
			return false, err
		}

		fmt.Println("Succesfully added relationship")
		return true, nil
	}

	//insert the new relationship
	statement, err = database.Db.Prepare("insert into follower(user_id,follower_id) values(?,?)")
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer statement.Close()
	_, err = statement.Exec(user.ID, followerID)

	if err != nil {
		log.Println(err)
		return false, err
	}

	fmt.Println("Succesfully added relationship")
	return true, nil
}

func (user *User) RemoveFollower(ctx context.Context, followerID int) (bool, error) {
	statement, err := database.Db.Prepare("update follower set active = false where user_id = ? and follower_id = ?")

	if err != nil {
		log.Println(err)
		return false, err
	}

	defer statement.Close()
	res, err := statement.Exec(user.ID, followerID)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false, err
	}

	count, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		return false, err
	}

	if count == 0 {
		fmt.Println("No relationship exist")
		err := fmt.Errorf("No relationship exist")
		return false, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Println("Succefully removed the relationship with ID", lastInsertID)
	return true, nil
}

func (user *User) FetchFeed(ctx context.Context) ([]*model.Tweet, error) {
	var tweets []*model.Tweet

	//fetch all the tweets from its following along with its own tweets
	statement, err := database.Db.Prepare("select t.t_id, t.u_id, t.created_on, t.content from tweet t LEFT JOIN follower f ON t.u_id = f.user_id where (f.active = true OR f.active is NULL) AND (f.follower_id = ? OR t.u_id = ?) ORDER BY t.created_on")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer statement.Close()
	rows, err := statement.QueryContext(ctx, user.ID, user.ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	hasNext := rows.Next()

	for ; hasNext; hasNext = rows.Next() {
		var tweet model.Tweet
		if err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.CreatedOn, &tweet.Content); err != nil {
			log.Println(err)
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	if len(tweets) == 0 {
		log.Println("No tweets are available for the user's feed: ", user.ID)
		return nil, fmt.Errorf("No tweets are available for the user's feed: ")
	}
	return tweets, nil

}

func GetUserIdByEmail(email string) (int, error) {
	statement, err := database.Db.Prepare("select u_id from user WHERE email = ?")
	if err != nil {
		log.Println(err)
		return -1, err
	}
	defer statement.Close()
	row := statement.QueryRow(email)

	var ID int
	err = row.Scan(&ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return -1, err
	} else if err == sql.ErrNoRows {
		log.Println(err)
		return -1, err
	}
	return ID, nil
}
