package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// User represents a user in the system. Users are typically created via OAuth
// using the AuthService but users can also be created directly for testing.
type User struct {
	ID        int64       `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Email     null.String `json:"email"`
	Name      string      `json:"name"`
	Password  null.String `json:"-"`

	// List of associated OAuth authentication objects.
	// Currently only GitHub is supported so there should only be a maximum of one.
	Auths []*Auth `json:"auths"`
}

// Users is a list of users.
type Users []*User

// Validate returns an error if the user contains invalid fields.
// This only performs basic validation.
func (u *User) Validate() error {

	if u.Name == "" {
		return Errorf(EINVALID, "name is required")
	}

	if u.Email.Valid && u.Email.String == "" {
		return Errorf(EINVALID, "email cannot be empty if provided")
	}

	return nil
}

// AvatarURL returns a URL to the avatar image for the user.
// This loops over all auth providers to find the first available avatar.
// Currently only GitHub is supported. Returns blank string if no avatar URL available.
func (u *User) AvatarURL(size int) string {
	for _, auth := range u.Auths {
		if s := auth.AvatarURL(size); s != "" {
			return s
		}
	}
	return ""
}
