package users

import (
	"fmt"
	"strings"
	"sync"
)

type UserRepository interface {
	GetUser(name string) (*User, error)
	PutUser(username string)
}

type userRepository struct {
	users map[string]*User
	mu sync.RWMutex
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make(map[string]*User),
	}
}

// GetUser returns a user by name
func (ur *userRepository) GetUser(name string) (*User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	user, ok := ur.users[name]
	if !ok {
		return &User{}, fmt.Errorf("error getting user '%s': does not exist", name)
	}
	return user, nil
}

// PutUser adds or updates a user in the database
func (ur *userRepository) PutUser(username string) {
	displayName := strings.Split(username, "@")[0]
	user := NewUser(username, displayName)

	ur.mu.Lock()
	defer ur.mu.Unlock()
	ur.users[user.name] = user
}
