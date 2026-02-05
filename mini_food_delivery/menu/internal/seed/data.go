package seed

import (
	"time"
)

type SeedMenu struct {
	Name        string
	Description string
	Active      bool
	Categories  []SeedCategory
}

type SeedCategory struct {
	Name     string
	Position int32
	Items    []SeedItem
}

type SeedItem struct {
	Name        string
	Description string
	Available   bool
	Prices      []SeedPrice
}

type SeedPrice struct {
	PriceCents int32
	Currency   string
	ValidFrom  time.Time
	ValidTo    *time.Time
}

var Menus = []SeedMenu{
	{
		Name:        "Main Menu",
		Description: "Default restaurant menu",
		Active:      true,
		Categories: []SeedCategory{
			{
				Name:     "Burgers",
				Position: 1,
				Items: []SeedItem{
					{
						Name:        "Classic Burger",
						Description: "Beef patty, cheddar, lettuce, tomato",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 1299,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -2, 0),
							},
						},
					},
					{
						Name:        "Double Cheeseburger",
						Description: "Two beef patties, double cheddar",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 1599,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -1, 0),
							},
						},
					},
					{
						Name:        "Vegan Burger",
						Description: "Plant-based patty, vegan mayo",
						Available:   false,
						Prices: []SeedPrice{
							{
								PriceCents: 1499,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -3, 0),
							},
						},
					},
				},
			},
			{
				Name:     "Sides",
				Position: 2,
				Items: []SeedItem{
					{
						Name:        "French Fries",
						Description: "Crispy golden fries",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 399,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -6, 0),
							},
						},
					},
					{
						Name:        "Onion Rings",
						Description: "Beer-battered onion rings",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 499,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -4, 0),
							},
						},
					},
				},
			},
			{
				Name:     "Drinks",
				Position: 3,
				Items: []SeedItem{
					{
						Name:        "Cola",
						Description: "Chilled carbonated drink",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 199,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -12, 0),
							},
						},
					},
					{
						Name:        "Lemonade",
						Description: "Fresh lemon juice and sugar",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 249,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -8, 0),
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        "Breakfast Menu",
		Description: "Morning specials served until 11am",
		Active:      true,
		Categories: []SeedCategory{
			{
				Name:     "Breakfast",
				Position: 1,
				Items: []SeedItem{
					{
						Name:        "Pancakes",
						Description: "Stack of pancakes with maple syrup",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 899,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -5, 0),
							},
						},
					},
					{
						Name:        "Omelette",
						Description: "Three eggs, cheese, choice of fillings",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 999,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -5, 0),
							},
						},
					},
				},
			},
			{
				Name:     "Hot Drinks",
				Position: 2,
				Items: []SeedItem{
					{
						Name:        "Coffee",
						Description: "Freshly brewed coffee",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 249,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -12, 0),
							},
						},
					},
					{
						Name:        "Tea",
						Description: "Black or green tea",
						Available:   true,
						Prices: []SeedPrice{
							{
								PriceCents: 199,
								Currency:   "USD",
								ValidFrom:  time.Now().AddDate(0, -12, 0),
							},
						},
					},
				},
			},
		},
	},
}
