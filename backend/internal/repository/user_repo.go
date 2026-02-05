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

type FirestoreUserRepo struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{
		client:     client,
		collection: "users",
	}
}

func (r *FirestoreUserRepo) Save(ctx context.Context, user *model.User) error {
	_, err := r.client.Collection(r.collection).Doc(user.ID).Set(ctx, user)
	if err != nil {
		return fmt.Errorf("saving user: %w", err)
	}
	return nil
}

func (r *FirestoreUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("getting user: %w", err)
	}

	var user model.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("parsing user: %w", err)
	}
	return &user, nil
}

func (r *FirestoreUserRepo) GetByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	iter := r.client.Collection(r.collection).Where("github_id", "==", githubID).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("user not found with github_id: %d", githubID)
	}
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}

	var user model.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("parsing user: %w", err)
	}
	return &user, nil
}

func (r *FirestoreUserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}
