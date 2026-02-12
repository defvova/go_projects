package menuitemprice

import (
	"context"

	"github.com/defvova/go_projects/mini_food_delivery/menu/db"
)

type MenuItemPriceStore interface {
	GetAllMenuItemPrices(ctx context.Context, menuItemId int64) ([]db.MenuItemPrice, error)
}

type Store struct {
	q *db.Queries
}

func NewStore(q *db.Queries) *Store {
	return &Store{q: q}
}

func (s *Store) GetAllMenuItemPrices(ctx context.Context, menuItemId int64) ([]db.MenuItemPrice, error) {
	data, err := s.q.GetMenuItemPrices(ctx, menuItemId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
