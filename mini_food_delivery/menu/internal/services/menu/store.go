package menu

import (
	"context"
	"mini_food_delivery/menu/db"
)

type MenuStore interface {
	GetAllMenus(ctx context.Context) ([]db.Menu, error)
}

type Store struct {
	q *db.Queries
}

func NewStore(q *db.Queries) *Store {
	return &Store{q: q}
}

func (s *Store) GetAllMenus(ctx context.Context) ([]db.Menu, error) {
	data, err := s.q.GetMenus(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}
