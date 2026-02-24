package auth

import (
	"context"
	"errors"

	repo "github.com/Ajay01103/goTransactonsAPI/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type postgresAuthRepository struct {
	queries *repo.Queries
}

// NewPostgresRepository constructs an auth Repository backed by sqlc-generated Queries.
// All sqlc and pgtype details are contained within this file — nothing leaks outward.
func NewPostgresRepository(queries *repo.Queries) Repository {
	return &postgresAuthRepository{queries: queries}
}

func (r *postgresAuthRepository) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	var profilePic pgtype.Text
	if params.ProfilePicture != "" {
		profilePic = pgtype.Text{String: params.ProfilePicture, Valid: true}
	}

	row, err := r.queries.CreateUser(ctx, repo.CreateUserParams{
		ID:             params.ID,
		Name:           params.Name,
		Email:          params.Email,
		Password:       params.Password,
		ProfilePicture: profilePic,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailTaken
		}
		return User{}, err
	}

	return User{
		ID:             row.ID,
		Name:           row.Name,
		Email:          row.Email,
		Password:       row.Password,
		ProfilePicture: row.ProfilePicture.String,
		CreatedAt:      row.CreatedAt.Time.String(),
	}, nil
}

func (r *postgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:             row.ID,
		Name:           row.Name,
		Email:          row.Email,
		Password:       row.Password,
		ProfilePicture: row.ProfilePicture.String,
		CreatedAt:      row.CreatedAt.Time.String(),
	}, nil
}
