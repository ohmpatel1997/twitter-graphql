package tweets

import (
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
	"log"
	"time"
)

type Tweet struct {
	UserID    int       `json:"u_id"`
	Content   string    `json:"content"`
	CreatedOn time.Time `josn:"created_on"`
}

func (tweet *Tweet) Save() (int64, error) {

	statement, err := database.Db.Prepare("INSERT INTO tweet(u_id,content,created_on) VALUES(?,?,?)")
	if err != nil {
		log.Println(err)
		return -1, err
	}

	res, err := statement.Exec(tweet.UserID, tweet.Content, tweet.CreatedOn)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	log.Print("Row inserted!")
	return id, nil
}
