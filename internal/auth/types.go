package auth

import "context"

// ── Domain model ─────────────────────────────────────────────────────────────

// User is the internal domain model — contains no storage-layer types.
type User struct {
	ID             string
	Name           string
	Email          string
	Password       string // bcrypt hash
	ProfilePicture string
	CreatedAt      string
}

// ── Service DTOs ──────────────────────────────────────────────────────────────

// RegisterInput is the DTO passed from handler → service for registration.
type RegisterInput struct {
	Name           string
	Email          string
	Password       string
	ProfilePicture string
}

// LoginInput is the DTO passed from handler → service for login.
type LoginInput struct {
	Email    string
	Password string
}

// UserPayload is the public user object embedded in auth responses.
type UserPayload struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// AuthResponse is returned after a successful register or login.
type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
}

// ── Repository DTO ────────────────────────────────────────────────────────────

// CreateUserParams carries the data needed to persist a new user.
// Uses only plain Go types — no sqlc or pgtype here.
type CreateUserParams struct {
	ID             string
	Name           string
	Email          string
	Password       string // already hashed
	ProfilePicture string
}

// ── Contracts ─────────────────────────────────────────────────────────────────

// Repository defines the data-access contract for the auth domain.
// All method signatures use domain types only.
type Repository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

// Service defines the business-logic contract for the auth domain.
type Service interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}
