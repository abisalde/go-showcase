package main

import (
	"context"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/abisalde/go-showcase/server/proto/church"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InMemoryChurch = make(map[string]*church.Church)
	mu             sync.Mutex
)

var (
	nonWordSpaceHyphen = regexp.MustCompile(`[^\w\s-]+`)
	whitespace         = regexp.MustCompile(`\s+`)
	multiHyphen        = regexp.MustCompile(`-{2,}`)
)

func generateChurchID(name string) string {
	mu.Lock()
	defer mu.Unlock()
	s := strings.ToLower(name)
	s = nonWordSpaceHyphen.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	s = whitespace.ReplaceAllString(s, "-")
	s = multiHyphen.ReplaceAllString(s, "-")

	log.Printf("ID was generated::: %v \n", s)
	return s
}

func (*Server) CreateChurch(ctx context.Context, in *church.Church) (*church.Church, error) {
	log.Printf("Create Church was invoked %v \n", in)

	id := generateChurchID(in.Name)

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

func (*Server) ListChurches(ctx context.Context, in *church.ListChurchesRequest) (*church.ListChurchesResponse, error) {
	log.Println("List Churches was invoked")

	mu.Lock()
	defer mu.Unlock()

	var churches []*church.Church
	for _, c := range InMemoryChurch {
		churches = append(churches, c)
	}

	response := &church.ListChurchesResponse{
		Churches: churches,
	}

	return response, nil
}
