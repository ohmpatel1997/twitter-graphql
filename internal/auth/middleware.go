package auth

import (
	"context"

	"log"

	"github.com/ohmpatel1997/twitter-graphql/internal/pkg/jwt"
	"github.com/ohmpatel1997/twitter-graphql/internal/users"
	"net/http"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	email string
}

func Middleware(next http.Handler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if len(header) == 0 {
			log.Println("Please log in first")
			http.Error(w, "Access Denied", http.StatusForbidden)
			return
		}

		tokenStr := header
		email, err := jwt.ParseToken(tokenStr)
		if err != nil {
			log.Println("invalid token")
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		user := users.User{
			Email: email,
		}

		id, err := users.GetUserIdByEmail(email)
		if err != nil {
			log.Println("User not found")
			http.Error(w, "Unauthorized user.", http.StatusUnauthorized)
			return
		}

		user.ID = id
		// put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, &user)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
