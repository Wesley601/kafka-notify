package db

import (
	"errors"
	"sync"

	"github.com/wesley601/kafka-notify/models"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

type UserNotifications map[string][]models.Notification

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func NewNotificationStore() *NotificationStore {
	return &NotificationStore{
		data: make(UserNotifications),
	}
}

func (ns *NotificationStore) Add(userID string, notification models.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) ([]models.Notification, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	notifications, ok := ns.data[userID]
	if !ok {
		return nil, ErrNoUserFound
	}

	return notifications, nil
}
