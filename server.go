package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/zrwaite/OweYeah/database"
	"github.com/zrwaite/OweYeah/graph"
	"github.com/zrwaite/OweYeah/graph/generated"
	"github.com/zrwaite/OweYeah/settings"
)

const defaultPort = "8080"

func main() {
	godotenv.Load(".env")
	settings.MatchDev()
	database.ConnectToMongoDB()
	// database.InitializeDatabase()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
