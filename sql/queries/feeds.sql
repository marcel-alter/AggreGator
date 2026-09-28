-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT 
    *,
    users.name AS owner

FROM 
    feeds
    LEFT JOIN users ON feeds.user_id = users.id
WHERE
    feeds.name = $1
    OR
    feeds.url = $1;

-- name: GetAllFeeds :many
SELECT 
    *,
    users.name AS owner
FROM 
    feeds
    LEFT JOIN users ON feeds.user_id = users.id;
