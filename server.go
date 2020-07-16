package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ohmpatel1997/twitter-graphql/graph"
	"github.com/ohmpatel1997/twitter-graphql/graph/generated"
	"github.com/ohmpatel1997/twitter-graphql/internal/auth"
	database "github.com/ohmpatel1997/twitter-graphql/internal/pkg/db/mysql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	http.HandleFunc("/query", auth.Middleware(srv))
	http.Handle("/query/users", srv)
	http.Handle("/mutation", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
