package main

import (
	"GRAPHQL/graph"
	"GRAPHQL/postgres"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err in load .env file:%s", err)
	}
	db, err := postgres.ConnTODB()
	if err != nil {
		log.Fatalf("failed to conn to db:%s", err)
	}
	resolver := graph.NewResolver(db)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)
	log.Printf("connect to http://localhost:%v/ for GraphQL playground", 8082)
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatalf("err in start serv:%s", err)
	}
}
