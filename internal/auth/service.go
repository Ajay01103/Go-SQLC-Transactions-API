package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	repo      Repository
	jwtSecret string
}

// NewService wires an auth Repository and JWT secret into a Service.
func NewService(repo Repository, jwtSecret string) Service {
	return &svc{repo: repo, jwtSecret: jwtSecret}
}

// Register hashes the password then delegates persistence to the repository.
func (s *svc) Register(ctx context.Context, input RegisterInput) (AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, CreateUserParams{
		ID:             cuid.New(),
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashed),
		ProfilePicture: input.ProfilePicture,
	})
	if err != nil {
		// propagate sentinel errors (e.g. ErrEmailTaken) unwrapped so callers can detect them
		if errors.Is(err, ErrEmailTaken) {
			return AuthResponse{}, err
		}
		return AuthResponse{}, fmt.Errorf("creating user: %w", err)
	}

	token, err := generateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("generating token: %w", err)
	}

	return AuthResponse{
		AccessToken: token,
		User: UserPayload{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			CreatedAt:      user.CreatedAt,
		},
	}, nil
}

// Login looks up the user by email and verifies the bcrypt password.
func (s *svc) Login(ctx context.Context, input LoginInput) (AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := generateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("generating token: %w", err)
	}

	return AuthResponse{
		AccessToken: token,
		User: UserPayload{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			CreatedAt:      user.CreatedAt,
		},
	}, nil
}
