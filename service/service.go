package service

import (
	"GRAPHQL/graph/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Register(ctx context.Context, input model.RegisterInput) (*model.User, error)
	Login(ctx context.Context, input model.LoginInput) (*model.Login, error)
	Logout(ctx context.Context, id string) (bool, error)
	AddCar(ctx context.Context, input model.CarInput) (*model.Car, error)
	UpdateCar(ctx context.Context, id string, input model.CarInput) (*model.Car, error)
	DeleteCar(ctx context.Context, id string) (*bool, error)
	GetAllCars(ctx context.Context) ([]*model.Car, error)
	GetCarByID(ctx context.Context, id string) (*model.Car, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	CarPublished(ctx context.Context) (<-chan *model.Car, error)
	ValidateSession(sessionID string) (bool, error)
	//AuthMiddleware(next http.Handler) http.Handler
}

type postgresService struct {
	db *pgxpool.Pool
}

func NewPostgresService(db *pgxpool.Pool) Service {
	return &postgresService{
		db: db,
	}
}
