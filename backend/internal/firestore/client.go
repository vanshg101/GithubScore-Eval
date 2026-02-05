package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewClient(projectID, credentialsFile string) *firestore.Client {
	ctx := context.Background()

	var client *firestore.Client
	var err error

	if credentialsFile != "" {
		client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	} else {
		client, err = firestore.NewClient(ctx, projectID)
	}

	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	return client
}
