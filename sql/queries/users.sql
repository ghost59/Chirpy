-- name: CreateUser :one

INSERT INTO users (id, created_at, updated_at, email, password, is_chirpy_red)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)

RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UserLookup :one
SELECT *
FROM users
WHERE email = $1;
-- name: UpdateUser :one
UPDATE users SET password = $1, email = $2
WHERE id  =  $3
RETURNING email, id, created_at, updated_at;

-- name: UpgradeUser :one
UPDATE users SET is_chirpy_red = true
WHERE id = $1
RETURNING *;
-- name: IDLookUp :one
SELECT *
FROM users
WHERE id = $1;
