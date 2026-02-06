package menuitem

import (
	"context"
	"mini_food_delivery/menu/db"
)

type MenuItemStore interface {
	GetAllMenuItems(ctx context.Context, categoryId int64) ([]db.MenuItem, error)
}

type Store struct {
	q *db.Queries
}

func NewStore(q *db.Queries) *Store {
	return &Store{q: q}
}

func (s *Store) GetAllMenuItems(ctx context.Context, categoryId int64) ([]db.MenuItem, error) {
	data, err := s.q.GetMenuItems(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
