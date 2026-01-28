-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users set email = $2, password_hash = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetSession :one
SELECT user_id, expires_at FROM sessions
WHERE token_hash = $1 LIMIT 1;

-- name: CreateSession :exec
INSERT INTO sessions (
    expires_at, user_id, token_hash
) VALUES (
    $1, $2, $3
);

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token_hash = $1;

-- name: GetRedirects :many
SELECT * FROM redirects
WHERE user_id = $1 ORDER BY created_at;

-- name: GetRedirectByToken :one
SELECT url FROM redirects
WHERE token = $1 LIMIT 1;

-- name: CreateRedirect :one
INSERT INTO redirects (
    url, user_id, token
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: DeleteRedirect :exec
DELETE FROM redirects
WHERE id = $1 AND user_id = $2;
