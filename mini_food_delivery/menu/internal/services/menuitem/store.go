package menuitem

import (
	"context"

	"github.com/defvova/go_projects/mini_food_delivery/menu/db"
)

type MenuItemStore interface {
	GetAllMenuItemsWithPrice(ctx context.Context, categoryId int64) ([]db.GetMenuItemsWithPriceByCategoryIdRow, error)
}

type Store struct {
	q *db.Queries
}

func NewStore(q *db.Queries) *Store {
	return &Store{q: q}
}

func (s *Store) GetAllMenuItemsWithPrice(ctx context.Context, categoryId int64) ([]db.GetMenuItemsWithPriceByCategoryIdRow, error) {
	data, err := s.q.GetMenuItemsWithPriceByCategoryId(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
