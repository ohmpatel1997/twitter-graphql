package tweets

import (
	"context"
	"github.com/ohmpatel1997/twitter-graphql/graph/model"
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
	"log"
	"time"
)

type Tweet struct {
	TweetID   int       `json:"t_id"`
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
	log.Print("Tweet inserted!")
	return id, nil
}

func (tweet *Tweet) FetchTweet(ctx context.Context) ([]*model.Tweet, error) {
	statement, err := database.Db.Prepare("select t_id, u_id, TIMESTAMP(created_on), content from tweet WHERE t_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := statement.QueryContext(ctx, tweet.TweetID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	hasNext := rows.Next()
	var tweets []*model.Tweet
	for ; hasNext; hasNext = rows.Next() {
		var tweet model.Tweet
		if err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.CreatedOn, &tweet.Content); err != nil {
			log.Println(err)
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}
