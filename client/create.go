package main

import (
	"context"
	"log"

	"github.com/abisalde/go-showcase/server/proto/church"
)

func create(client church.ChurchServiceClient) {
	log.Println("Creating a new church...")

	newChurch := &church.Church{
		Name:    "Lane Book Church",
		Address: "123 Main St, Anytown, USA",
		Pastor:  "Mark Batterson",
	}

	resp, err := client.CreateChurch(context.Background(), newChurch)
	if err != nil {
		log.Fatalf("Error while creating church: %v", err)
	}

	log.Printf("Church created: %v", resp.Id)
}

func getChurch(client church.ChurchServiceClient) {
	log.Println("Getting church details...")

	req := &church.GetChurchRequest{
		Id: "1",
	}

	resp, err := client.GetChurch(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while getting church: %v", err)
	}

	log.Printf("Church details: %v", resp.Church)
}
