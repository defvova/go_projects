-- name: GetMenu :one
SELECT * FROM menus
WHERE id = $1 LIMIT 1;

-- name: GetMenus :many
SELECT * FROM menus;
