package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/ohmpatel1997/twitter-graphql/graph/generated"
	"github.com/ohmpatel1997/twitter-graphql/graph/model"
	"github.com/ohmpatel1997/twitter-graphql/internal/tweets"
	"github.com/ohmpatel1997/twitter-graphql/internal/users"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, input model.NewTweet) (*model.Tweet, error) {
	var tweet tweets.Tweet
	userIntID, err := strconv.Atoi(input.UserID)

	if err != nil {
		return nil, err
	}

	tweet.UserID = userIntID
	tweet.Content = input.Content
	tweet.CreatedOn = time.Now()
	ID, err := tweet.Save()

	if err != nil {
		return nil, err
	}

	return &model.Tweet{
		TweetID:   strconv.FormatInt(ID, 10),
		CreatedOn: tweet.CreatedOn,
		UserID:    strconv.FormatInt(int64(tweet.UserID), 10),
		Content:   tweet.Content,
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	var user users.User
	user.FirstName = input.FirstName

	if len(*input.LastName) > 0 {
		user.LastName = *input.LastName
	}
	user.Username = input.Username
	hashedPass, err := HashPassword(input.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	user.Password = hashedPass
	user.Email = input.Email
	userID, err := user.Save()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	user.ID = int(userID)

	return &model.User{
		UserID:    strconv.FormatInt(int64(user.ID), 10),
		FirstName: user.FirstName,
		LastName:  &user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Deleted:   user.Deleted,
	}, nil

}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (bool, error) {
	var user users.User
	user.Username = *input.Username
	user.Email = *input.Email
	user.Password = input.Password
	correct, err := user.Authenticate()

	if !correct {
		invalidUser := &users.InvalidUsernameOrPasswordError{
			UserName: user.Username,
			Email:    user.Email,
		}
		log.Println(invalidUser.Error())
		return false, invalidUser
	}

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil

}

func (r *mutationResolver) CreateRelationship(ctx context.Context, input model.Relationship) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveRelationship(ctx context.Context, intput model.Relationship) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tweets(ctx context.Context) ([]*model.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
