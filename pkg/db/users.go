package db

import (
	"errors"
	"sync"

	"github.com/wesley601/kafka-notify/models"
)

var (
	ErrUserNotFoundInProducer = errors.New("user not found")
)

type UserStore struct {
	data []models.User
	mu   sync.RWMutex
}

func NewUserStore(users []models.User) *UserStore {
	return &UserStore{
		data: users,
	}
}

func (us *UserStore) Add(userID string, user models.User) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.data = append(us.data, user)
}

func (us *UserStore) ByID(userID int) (models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	for _, user := range us.data {
		if user.ID == userID {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFoundInProducer
}
