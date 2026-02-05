package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreRankingRepo struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreRankingRepo(client *firestore.Client) *FirestoreRankingRepo {
	return &FirestoreRankingRepo{
		client:     client,
		collection: "rankings",
	}
}

func (r *FirestoreRankingRepo) Save(ctx context.Context, ranking *model.Ranking) error {
	_, err := r.client.Collection(r.collection).Doc(ranking.SnapshotDate).Set(ctx, ranking)
	if err != nil {
		return fmt.Errorf("saving ranking: %w", err)
	}
	return nil
}

func (r *FirestoreRankingRepo) GetLatest(ctx context.Context) (*model.Ranking, error) {
	iter := r.client.Collection(r.collection).OrderBy("created_at", firestore.Desc).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("getting latest ranking: %w", err)
	}

	var ranking model.Ranking
	if err := doc.DataTo(&ranking); err != nil {
		return nil, fmt.Errorf("parsing ranking: %w", err)
	}
	return &ranking, nil
}

func (r *FirestoreRankingRepo) GetByDate(ctx context.Context, date string) (*model.Ranking, error) {
	doc, err := r.client.Collection(r.collection).Doc(date).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("ranking not found for date: %s", date)
		}
		return nil, fmt.Errorf("getting ranking: %w", err)
	}

	var ranking model.Ranking
	if err := doc.DataTo(&ranking); err != nil {
		return nil, fmt.Errorf("parsing ranking: %w", err)
	}
	return &ranking, nil
}
