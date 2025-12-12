package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/abisalde/go-showcase/server/pkg/middleware"
	"github.com/abisalde/go-showcase/server/proto/church"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var greetWithDeadlineTime time.Duration = 1 * time.Second

type Server struct {
	church.UnimplementedChurchServiceServer
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}
	addr := "0.0.0.0:" + port

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	defer lis.Close()
	log.Printf("Listening at %s\n", addr)

	opts := []grpc.ServerOption{}

	tlsEnabled := os.Getenv("TLS_ENABLED")
	tls := tlsEnabled == "true"

	log.Printf("Is TLS Enabled::: %v", tls)

	if tls {
		certFile := "./ssl/server.crt"
		keyFile := "./ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v\n", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(middleware.LogInterceptor(), middleware.CheckHeaderInterceptor()))

	s := grpc.NewServer(opts...)
	church.RegisterChurchServiceServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
