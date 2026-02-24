package users

import (
	"context"

	repo "github.com/Ajay01103/goTransactonsAPI/internal/adapters/postgresql/sqlc"
)

type postgresRepository struct {
	queries *repo.Queries
}

// NewPostgresRepository constructs a users Repository backed by sqlc-generated Queries.
// All sqlc and pgtype details are contained within this file — nothing leaks outward.
func NewPostgresRepository(queries *repo.Queries) Repository {
	return &postgresRepository{queries: queries}
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id string) (UserRecord, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return UserRecord{}, err
	}

	return UserRecord{
		ID:             row.ID,
		Name:           row.Name,
		Email:          row.Email,
		ProfilePicture: row.ProfilePicture.String,
		CreatedAt:      row.CreatedAt.Time.String(),
	}, nil
}
