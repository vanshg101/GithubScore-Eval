package store

import (
	"fmt"
	"sync"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

type UserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*model.User),
	}
}

func (s *UserStore) Save(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return nil
}

func (s *UserStore) GetByID(id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return user, nil
}

func (s *UserStore) GetByGitHubID(githubID int64) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.GitHubID == githubID {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found with github_id: %d", githubID)
}

func (s *UserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, id)
	return nil
}
