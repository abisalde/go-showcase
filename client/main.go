package main

import (
	"log"
	"os"

	"github.com/abisalde/go-showcase/server/proto/church"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := os.Getenv("GRPC_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:50052"
	}

	tlsEnabled := os.Getenv("TLS_ENABLED")
	tls := tlsEnabled == "true"
	opts := []grpc.DialOption{}

	if tls {
		certFile := "./ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		log.Println("Cert File", sslErr)
		if sslErr != nil {
			log.Fatalf("Failed loading CA trust certificate: %v\n", sslErr)
			return
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
	defer conn.Close()

	c := church.NewChurchServiceClient(conn)

	// create(c)
	getChurch(c)

	log.Println("✅ Client connected", c)
}
