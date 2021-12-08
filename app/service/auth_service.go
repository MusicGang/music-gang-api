package service

import (
	"context"

	"github.com/music-gang/music-gang-api/app/entity"
)

// AuthService is an interface for managing user authentication
type AuthService interface {
	// CreateAuth creates a new auth.
	// If is attached to a user, links the auth to the user, otherwise creates a new user.
	// On success, the auth.ID is set.
	CreateAuth(ctx context.Context, auth *entity.Auth) error

	// DeleteAuth deletes an auth.
	// Do not delete underlying user.
	DeleteAuth(ctx context.Context, id int64) error

	// FindAuthByID returns a single auth by its id.
	FindAuthByID(ctx context.Context, id int64) (*entity.Auth, error)

	// FindAuths returns a list of auths.
	// Predicate can be used to filter the results.
	// Also returns the total count of auths.
	FindAuths(ctx context.Context, filter AuthFilter) (entity.Auths, int, error)
}

// AuthFilter represents a filter for auths.
type AuthFilter struct {
	ID       *int64  `json:"id"`
	UserID   *int64  `json:"user_id"`
	Source   *string `json:"source"`
	SourceID *string `json:"source_id"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
