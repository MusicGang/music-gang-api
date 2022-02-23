package entity

import (
	"time"

	"github.com/music-gang/music-gang-api/app/apperr"
)

// Visibility consts for the visibility of the contract.
const (
	VisibilityPrivate = "private"
	VisibilityPublic  = "public"
)

// Visibility defines the visibility of a contract.
type Visibility string

// Validate validates the visibility.
func (v Visibility) Validate() error {
	switch v {
	case
		VisibilityPrivate,
		VisibilityPublic:
		return nil
	default:
		return apperr.Errorf(apperr.EINVALID, "invalid visibility")
	}
}

// Contracts represents a list of contracts.
type Contracts []*Contract

// Contract represents a contract.
// The contract is a cloud function that is executed on a server, deployed by users;
// The contract can have multiple revisions.
type Contract struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	UserID      int64      `json:"user_id"`
	Visibility  Visibility `json:"visibility"`
	MaxFuel     Fuel       `json:"max_fuel"` // The maximum amount of fuel that can be burned from the contract.
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	LastRevision *Revision `json:"last_revision"`
	User         *User     `json:"user"`
}

// MaxExecutionTime returns the maximum execution time of the contract.
// MaxExecutionTime is based on max fuel compared with fuelAmountTable.
func (c *Contract) MaxExecutionTime() time.Duration {
	for t, fuel := range fuelAmountTable {
		if c.MaxFuel <= fuel {
			return t
		}
	}

	return MaxExecutionTime
}

// Validate validates the contract.
func (c *Contract) Validate() error {

	if c.Name == "" {
		return apperr.Errorf(apperr.EINVALID, "contract name is required")
	}

	if c.UserID == 0 {
		return apperr.Errorf(apperr.EINVALID, "User ID cannot be empty if provided")
	}

	if c.MaxFuel == 0 {
		return apperr.Errorf(apperr.EINVALID, "Max fuel is required")
	}

	if err := c.Visibility.Validate(); err != nil {
		return err
	}

	return nil
}
