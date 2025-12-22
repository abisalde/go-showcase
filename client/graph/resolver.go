package graph

import "github.com/abisalde/go-showcase/server/proto/church"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Client struct {
	churchClient church.ChurchServiceClient
}

type Resolver struct {
	*Client
}

func NewResolver(churchClient church.ChurchServiceClient) *Resolver {
	return &Resolver{
		Client: &Client{
			churchClient: churchClient,
		},
	}
}
