package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/tobi/contracts/backend/internal/model"
)

type Store interface {
	CreateUser(ctx context.Context, u model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, u model.User) error
	ListUsers(ctx context.Context) ([]model.User, error)

	GetSettings(ctx context.Context, userID string) (model.UserSettings, error)
	UpdateSettings(ctx context.Context, userID string, s model.UserSettings) error

	ListCategories(ctx context.Context, userID string) ([]model.Category, error)
	GetCategory(ctx context.Context, userID string, id uuid.UUID) (model.Category, error)
	CreateCategory(ctx context.Context, userID string, c model.Category) error
	UpdateCategory(ctx context.Context, userID string, c model.Category) error
	DeleteCategory(ctx context.Context, userID string, id uuid.UUID) error

	ListContracts(ctx context.Context, userID string) ([]model.Contract, error)
	ListContractsByCategory(ctx context.Context, userID string, categoryID uuid.UUID) ([]model.Contract, error)
	GetContract(ctx context.Context, userID string, id uuid.UUID) (model.Contract, error)
	CreateContract(ctx context.Context, userID string, c model.Contract) error
	UpdateContract(ctx context.Context, userID string, c model.Contract) error
	DeleteContract(ctx context.Context, userID string, id uuid.UUID) error

	Close() error
}
