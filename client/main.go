package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abisalde/go-showcase/client/graph"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("CLIENT_PORT")
	if port == "" {
		port = defaultPort
	}

	churchClient := SetupGRPCServer()

	resolvers := graph.NewResolver(churchClient)
	srv := GraphQLServer(resolvers)

	log.Println("✅ Client connected", churchClient)

	http.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
