package repository

import (
	"errors"
	"sync"

	"github.com/cushydigit/microstore/shared/types"
)

type InMemoryUserRepo struct {
	users map[string]*types.User
	mu    sync.Mutex
	idSeq int
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*types.User),
		idSeq: 1,
	}
}

func (r *InMemoryUserRepo) FindByEmail(email string) (*types.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepo) Create(user *types.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	user.ID = r.idSeq
	r.idSeq++
	r.users[user.Email] = user
	return nil
}
