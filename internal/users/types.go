package users

import "context"

// ── Domain model ──────────────────────────────────────────────────────────────

// UserRecord is the internal domain model — no storage-layer types.
type UserRecord struct {
	ID             string
	Name           string
	Email          string
	ProfilePicture string
	CreatedAt      string
}

// ── Service DTO ───────────────────────────────────────────────────────────────

// UserResponse is the public DTO returned from service → handler.
type UserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// ── Contracts ─────────────────────────────────────────────────────────────────

// Repository defines the data-access contract for the users domain.
// All method signatures use domain types only — no sqlc or pgtype.
type Repository interface {
	GetUserByID(ctx context.Context, id string) (UserRecord, error)
}

// Service defines the business-logic contract for the users domain.
type Service interface {
	GetCurrentUser(ctx context.Context, userID string) (UserResponse, error)
}