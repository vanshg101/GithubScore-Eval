package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreDeveloperRepo struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreDeveloperRepo(client *firestore.Client) *FirestoreDeveloperRepo {
	return &FirestoreDeveloperRepo{
		client:     client,
		collection: "developers",
	}
}

func (r *FirestoreDeveloperRepo) Save(ctx context.Context, dev *model.Developer) error {
	_, err := r.client.Collection(r.collection).Doc(dev.Username).Set(ctx, dev)
	if err != nil {
		return fmt.Errorf("saving developer: %w", err)
	}
	return nil
}

func (r *FirestoreDeveloperRepo) GetByUsername(ctx context.Context, username string) (*model.Developer, error) {
	doc, err := r.client.Collection(r.collection).Doc(username).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("developer not found: %s", username)
		}
		return nil, fmt.Errorf("getting developer: %w", err)
	}

	var dev model.Developer
	if err := doc.DataTo(&dev); err != nil {
		return nil, fmt.Errorf("parsing developer: %w", err)
	}
	return &dev, nil
}

func (r *FirestoreDeveloperRepo) GetAll(ctx context.Context) ([]*model.Developer, error) {
	iter := r.client.Collection(r.collection).Documents(ctx)
	defer iter.Stop()

	var developers []*model.Developer
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterating developers: %w", err)
		}

		var dev model.Developer
		if err := doc.DataTo(&dev); err != nil {
			return nil, fmt.Errorf("parsing developer: %w", err)
		}
		developers = append(developers, &dev)
	}
	return developers, nil
}

func (r *FirestoreDeveloperRepo) Delete(ctx context.Context, username string) error {
	_, err := r.client.Collection(r.collection).Doc(username).Delete(ctx)
	if err != nil {
		return fmt.Errorf("deleting developer: %w", err)
	}
	return nil
}
