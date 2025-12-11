package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/abisalde/go-showcase/server/proto/church"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InMemoryChurch  = make(map[string]*church.Church)
	churchIDCounter int64
	mu              sync.Mutex
)

func generateChurchID() string {
	mu.Lock()
	defer mu.Unlock()
	churchIDCounter++
	log.Printf("ID was generated::: %v \n", churchIDCounter)
	return fmt.Sprintf("%d", churchIDCounter)
}

func (*Server) CreateChurch(ctx context.Context, in *church.Church) (*church.Church, error) {
	log.Printf("Create Church was invoked %v \n", in)

	id := generateChurchID()

	newChurch := &church.Church{
		Id:      id,
		Name:    in.Name,
		Address: in.Address,
		Pastor:  in.Pastor,
	}

	mu.Lock()
	InMemoryChurch[id] = newChurch
	mu.Unlock()

	log.Printf("Church created with ID::: %s", id)

	return newChurch, nil
}

func (*Server) GetChurch(ctx context.Context, in *church.GetChurchRequest) (*church.GetChurchResponse, error) {

	if in.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Church ID is required")
	}

	mu.Lock()
	churchData, exists := InMemoryChurch[in.Id]
	mu.Unlock()

	if !exists {
		return nil, status.Error(codes.NotFound, "Church not found")
	}

	response := &church.GetChurchResponse{
		Church: churchData,
	}

	return response, nil
}
