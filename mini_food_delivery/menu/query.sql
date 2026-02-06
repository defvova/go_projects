-- name: GetMenu :one
SELECT * FROM menus
WHERE id = $1 LIMIT 1;

-- name: GetMenus :many
SELECT * FROM menus;

-- name: GetMenuByName :one
SELECT * FROM menus
WHERE name = $1 LIMIT 1;

-- name: CreateMenu :one
INSERT INTO menus (
    name, description, active
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCategoryByMenuAndName :one
SELECT * FROM categories
WHERE menu_id = $1 AND name = $2 LIMIT 1;

-- name: GetCategories :many
SELECT * FROM categories
WHERE menu_id = $1;

-- name: CreateCategory :one
INSERT INTO categories (
    menu_id, name, position
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMenuItemsByCategoryId :many
SELECT * FROM menu_items
WHERE category_id = $1;

-- name: GetMenuItemsWithPriceByCategoryId :many
SELECT * FROM menu_items
INNER JOIN menu_item_prices ON menu_item_prices.menu_item_id = menu_items.id
WHERE menu_items.category_id = $1;

-- name: GetMenuItemByCategoryAndName :one
SELECT * FROM menu_items
WHERE category_id = $1 AND name = $2 LIMIT 1;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
    category_id, name, description, available
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetMenuItemPrices :many
SELECT * FROM menu_item_prices
WHERE menu_item_id = $1;

-- name: PriceExists :one
SELECT EXISTS (
    SELECT 1 FROM menu_item_prices
    WHERE menu_item_id = $1
        AND price_cents = $2
        AND currency = $3
        AND valid_from = $4
);

-- name: CreateMenuItemPrice :one
INSERT INTO menu_item_prices (
    menu_item_id, price_cents, currency, valid_from, valid_to
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
