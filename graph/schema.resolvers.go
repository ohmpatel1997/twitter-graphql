package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ohmpatel1997/twitter-graphql/common"
	"github.com/ohmpatel1997/twitter-graphql/graph/generated"
	"github.com/ohmpatel1997/twitter-graphql/graph/model"
	"github.com/ohmpatel1997/twitter-graphql/internal/auth"
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
	"github.com/ohmpatel1997/twitter-graphql/internal/pkg/jwt"
	"github.com/ohmpatel1997/twitter-graphql/internal/tweets"
	"github.com/ohmpatel1997/twitter-graphql/internal/users"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, input model.NewTweet) (*model.Tweet, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		log.Println("can not able to fetch user from context")
		return nil, fmt.Errorf("Access denied")
	}

	var tweet tweets.Tweet

	tweet.UserID = user.ID //fetch tweet from current user
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

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.FirstName = input.FirstName

	if len(*input.LastName) > 0 {
		user.LastName = *input.LastName
	}
	user.Username = input.Username
	hashedPass, err := common.HashPassword(input.Password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	user.Password = hashedPass
	user.Email = input.Email

	userID, err := user.Save() //first save the user, auth token is not mandatory on sign up
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("Successfully user created with ID", userID)

	token, err := jwt.GenerateToken(user.Email)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
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
		return "", invalidUser
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	token, err := jwt.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateRelationship(ctx context.Context, input model.Relationship) (bool, error) {
	userFromContext := auth.ForContext(ctx)
	if userFromContext == nil {
		log.Println("can not able to fetch user from context")
		return false, fmt.Errorf("Access denied")
	}

	var user users.User
	var err error
	user.ID, err = strconv.Atoi(input.UserID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	FollowerID, err := strconv.Atoi(input.FollowerID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return user.AddFollower(ctx, FollowerID)
}

func (r *mutationResolver) RemoveRelationship(ctx context.Context, intput model.Relationship) (bool, error) {
	userFromContext := auth.ForContext(ctx)
	if userFromContext == nil {
		log.Println("can not able to fetch user from context")
		return false, fmt.Errorf("Access denied")
	}

	var user users.User
	var err error

	user.ID, err = strconv.Atoi(intput.UserID)
	if err != nil {
		log.Println(err)
		return false, err
	}

	FollowerID, err := strconv.Atoi(intput.FollowerID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return user.RemoveFollower(ctx, FollowerID)
}

func (r *queryResolver) Tweets(ctx context.Context, tweetID *string, userID *string, username *string) ([]*model.Tweet, error) {
	//querying based on userId or username
	if tweetID == nil && (username != nil || userID != nil) {

		user := users.User{}

		if username != nil {
			user.Username = *username
		}

		if userID != nil {
			if intID, err := strconv.Atoi(*userID); err == nil {
				user.ID = intID
			}
		}

		return user.FetchAllTweetsOfUser(ctx)
	}

	// fetch tweet using tweet_id
	tweet := tweets.Tweet{}
	if intID, err := strconv.Atoi(*tweetID); err == nil {
		tweet.TweetID = intID
	}
	return tweet.FetchTweet(ctx)
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	statement, err := database.Db.Prepare("select u_id,f_name,l_name,user_name,email from user where deleted = false")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := statement.QueryContext(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	hasNext := rows.Next()
	var users []*model.User
	for ; hasNext; hasNext = rows.Next() {
		user := model.User{}
		if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}

	if len(users) == 0 {
		log.Println("No users found")
		return nil, fmt.Errorf("No users found")
	}

	return users, nil
}

func (r *queryResolver) Feed(ctx context.Context) ([]*model.Tweet, error) {
	userFromContext := auth.ForContext(ctx)
	if userFromContext == nil {
		log.Println("can not able to fetch user from context")
		return nil, fmt.Errorf("Access denied")
	}

	return userFromContext.FetchFeed(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
