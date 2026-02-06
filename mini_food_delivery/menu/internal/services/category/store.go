package category

import (
	"context"
	"mini_food_delivery/menu/db"
)

type CategoryStore interface {
	GetAllCategories(ctx context.Context, menuId int64) ([]db.Category, error)
}

type Store struct {
	q *db.Queries
}

func NewStore(q *db.Queries) *Store {
	return &Store{q: q}
}

func (s *Store) GetAllCategories(ctx context.Context, menuId int64) ([]db.Category, error) {
	data, err := s.q.GetCategories(ctx, menuId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
