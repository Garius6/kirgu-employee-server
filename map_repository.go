package main

import (
	"errors"
	"sync"

	"kirgu.ru/employee/model"
)

type MapRepository struct {
	m     sync.Mutex
	users map[string]*model.User
}

func NewMapRepository() *MapRepository {
	return &MapRepository{users: make(map[string]*model.User), m: sync.Mutex{}}
}

func (m *MapRepository) SignIn(username string, password string) (*model.User, error) {
	if user, ok := m.users[username]; !ok {
		return nil, errors.New("user not found")
	} else {
		return user, nil
	}
}

func (m *MapRepository) SignUp(username string, password string, passwordConfirmation string) error {
	if _, ok := m.users[username]; ok {
		return errors.New("user alreadty exists")
	} else {
		m.users[username] = &model.User{Username: username, Password: password}
		return nil
	}
}
