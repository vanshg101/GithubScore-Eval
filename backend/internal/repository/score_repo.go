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

type FirestoreScoreRepo struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreScoreRepo(client *firestore.Client) *FirestoreScoreRepo {
	return &FirestoreScoreRepo{
		client:     client,
		collection: "scores",
	}
}

func (r *FirestoreScoreRepo) Save(ctx context.Context, score *model.Score) error {
	_, err := r.client.Collection(r.collection).Doc(score.Username).Set(ctx, score)
	if err != nil {
		return fmt.Errorf("saving score: %w", err)
	}
	return nil
}

func (r *FirestoreScoreRepo) GetByUsername(ctx context.Context, username string) (*model.Score, error) {
	doc, err := r.client.Collection(r.collection).Doc(username).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("score not found: %s", username)
		}
		return nil, fmt.Errorf("getting score: %w", err)
	}

	var score model.Score
	if err := doc.DataTo(&score); err != nil {
		return nil, fmt.Errorf("parsing score: %w", err)
	}
	return &score, nil
}

func (r *FirestoreScoreRepo) GetAll(ctx context.Context) ([]*model.Score, error) {
	iter := r.client.Collection(r.collection).Documents(ctx)
	defer iter.Stop()

	var scores []*model.Score
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterating scores: %w", err)
		}

		var score model.Score
		if err := doc.DataTo(&score); err != nil {
			return nil, fmt.Errorf("parsing score: %w", err)
		}
		scores = append(scores, &score)
	}
	return scores, nil
}

func (r *FirestoreScoreRepo) Delete(ctx context.Context, username string) error {
	_, err := r.client.Collection(r.collection).Doc(username).Delete(ctx)
	if err != nil {
		return fmt.Errorf("deleting score: %w", err)
	}
	return nil
}
