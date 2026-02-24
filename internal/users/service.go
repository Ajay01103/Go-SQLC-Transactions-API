package users

import (
	"context"
	"fmt"
)

type svc struct {
	repo Repository
}

// NewService wires a users Repository into a Service.
func NewService(repo Repository) Service {
	return &svc{repo: repo}
}

// GetCurrentUser fetches the user by ID from the repository.
func (s *svc) GetCurrentUser(ctx context.Context, userID string) (UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return UserResponse{}, fmt.Errorf("user not found")
	}

	return UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
	}, nil
}
