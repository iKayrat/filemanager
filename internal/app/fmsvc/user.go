package fmsvc

import (
	"context"
)

type userService interface {
	CreateUser(ctx context.Context, name string) (User, error)
	GetUser(ctx context.Context, id string) (User, error)
}

type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser implements userService.
func (ims *User) CreateUser(ctx context.Context, name string) (User, error) {

	return User{}, nil
}

// GetUser implements userService.
func (ims *User) GetUser(ctx context.Context, id string) (User, error) {

	return User{}, nil
}
