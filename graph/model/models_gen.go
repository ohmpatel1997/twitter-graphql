// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Login struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password string  `json:"password"`
}

type NewTweet struct {
	Content string `json:"content"`
}

type NewUser struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
}

type Relationship struct {
	UserID     string `json:"user_id"`
	FollowerID string `json:"follower_id"`
	Active     *bool  `json:"active"`
}

type Tweet struct {
	TweetID   string    `json:"tweet_id"`
	UserID    string    `json:"user_id"`
	CreatedOn time.Time `json:"created_on"`
	Content   string    `json:"content"`
}

type User struct {
	UserID    string  `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	Deleted   bool    `json:"deleted"`
}
