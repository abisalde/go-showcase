package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/abisalde/go-showcase/client/graph"
	"github.com/abisalde/go-showcase/server/proto/church"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GraphQLServer(resolver *graph.Resolver) *handler.Server {

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Client: resolver.Client,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func SetupGRPCServer() church.ChurchServiceClient {

	addr := os.Getenv("GRPC_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:50052"
	}

	tlsEnabled := os.Getenv("TLS_ENABLED")
	tls := tlsEnabled == "true"
	opts := []grpc.DialOption{}

	if tls {
		certFile := os.Getenv("TLS_CA_CERT")
		if certFile == "" {
			certFile = "./ssl/ca.crt" // Default for local development
		}
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		log.Printf("Loading certificate from: %s\n", certFile)
		if sslErr != nil {
			log.Fatalf("Failed loading CA trust certificate from %s: %v\n", certFile, sslErr)
			return nil
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
	}

	opts = append(opts, grpc.WithChainUnaryInterceptor(LogInterceptor(), AddHeaderInterceptor()))

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	c := church.NewChurchServiceClient(conn)

	return c
}
