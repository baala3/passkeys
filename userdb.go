package main

import (
	"fmt"
	"sync"
)

type UserDB struct {
	users map[string]*User
	mu sync.RWMutex
}

var db *UserDB

// DB returns a singleton instance of UserDB
func DB() *UserDB {
	if db == nil {
		db = &UserDB{
			users: make(map[string]*User),
		}
	}
	return db
}

// GetUser returns a user by name
func (userDB *UserDB) GetUser(name string) (*User, error) {
	userDB.mu.RLock()
	defer userDB.mu.RUnlock()
	user, ok := userDB.users[name]
	if !ok {
		return &User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

// PutUser adds or updates a user in the database
func (userDB *UserDB) PutUser(user *User) {
	userDB.mu.Lock()
	defer userDB.mu.Unlock()
	userDB.users[user.name] = user
}

// GetUserCount returns the number of users in the database
func (userDB *UserDB) GetUserCount() int {
	userDB.mu.RLock()
	defer userDB.mu.RUnlock()
	return len(userDB.users)
}

