package seed

import (
	"context"
	"mini_food_delivery/menu/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func SeedAll(ctx context.Context, q *db.Queries) error {
	for _, menu := range Menus {
		if err := seedMenu(ctx, q, menu); err != nil {
			return err
		}
	}
	return nil
}

func seedMenu(ctx context.Context, q *db.Queries, m SeedMenu) error {
	menu, err := q.GetMenuByName(ctx, m.Name)
	if err != nil {
		menu, err = q.CreateMenu(ctx, db.CreateMenuParams{
			Name: m.Name,
			Description: pgtype.Text{
				String: m.Description,
				Valid:  m.Description != "",
			},
			Active: m.Active,
		})
		if err != nil {
			return err
		}
	}

	for _, c := range m.Categories {
		if err := seedCategory(ctx, q, menu.ID, c); err != nil {
			return err
		}
	}

	return nil
}

func seedCategory(ctx context.Context, q *db.Queries, menuID int64, c SeedCategory) error {
	cat, err := q.GetCategoryByMenuAndName(ctx, db.GetCategoryByMenuAndNameParams{
		MenuID: menuID,
		Name:   c.Name,
	})
	if err != nil {
		cat, err = q.CreateCategory(ctx, db.CreateCategoryParams{
			MenuID:   menuID,
			Name:     c.Name,
			Position: c.Position,
		})
		if err != nil {
			return err
		}
	}

	for _, item := range c.Items {
		if err := seedItem(ctx, q, cat.ID, item); err != nil {
			return err
		}
	}

	return nil
}

func seedItem(ctx context.Context, q *db.Queries, categoryID int64, i SeedItem) error {
	item, err := q.GetMenuItemByCategoryAndName(ctx, db.GetMenuItemByCategoryAndNameParams{
		CategoryID: categoryID,
		Name:       i.Name,
	})
	if err != nil {
		item, err = q.CreateMenuItem(ctx, db.CreateMenuItemParams{
			CategoryID: categoryID,
			Name:       i.Name,
			Description: pgtype.Text{
				String: i.Description,
				Valid:  i.Description != "",
			},
			Available: i.Available,
		})
		if err != nil {
			return err
		}
	}

	for _, p := range i.Prices {
		if err := seedPrice(ctx, q, item.ID, p); err != nil {
			return err
		}
	}

	return nil
}

func seedPrice(ctx context.Context, q *db.Queries, menuItemID int64, p SeedPrice) error {
	exists, err := q.PriceExists(ctx, db.PriceExistsParams{
		MenuItemID: menuItemID,
		PriceCents: p.PriceCents,
		Currency:   p.Currency,
		ValidFrom:  p.ValidFrom,
	})
	if err != nil || exists {
		return err
	}

	_, err = q.CreateMenuItemPrice(ctx, db.CreateMenuItemPriceParams{
		MenuItemID: menuItemID,
		PriceCents: p.PriceCents,
		Currency:   p.Currency,
		ValidFrom:  p.ValidFrom,
		ValidTo:    p.ValidTo,
	})
	return err
}
