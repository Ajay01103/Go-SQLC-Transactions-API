package auth

import "errors"

// Sentinel errors for the auth domain.
var (
	// ErrEmailTaken is returned when a registration attempt uses an email that already exists.
	ErrEmailTaken = errors.New("email already in use")
)
