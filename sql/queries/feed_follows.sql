-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
    RETURNING *
    )
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    followers.name AS user_name,
    creators.name AS creator_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users AS followers ON inserted_feed_follow.user_id = followers.id
INNER JOIN users AS creators ON feeds.user_id = creators.id;

-- name: GetFeedFollow :one
SELECT 
    *,
    feeds.name AS feed_name,
    followers.name AS user_name,
    creators.name AS creator_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users AS followers ON feed_follows.user_id = followers.id
INNER JOIN users AS creators ON feeds.user_id = creators.id
WHERE
    feeds.url = $1;

-- name: GetFeedFollowForUser :many
SELECT
    *,
    feeds.name AS feed_name,
    followers.name AS user_name,
    creators.name AS creator_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users AS followers ON feed_follows.user_id = followers.id
INNER JOIN users AS creators ON feeds.user_id = creators.id
WHERE
    followers.name = $1;
