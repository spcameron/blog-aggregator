-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: GetPostsForUser :many
SELECT
    p.title,
    p.url,
    p.description,
    p.published_at,
    f.name AS feed_name,
    f.url AS feed_url
FROM
    posts AS p
    JOIN feeds AS f ON p.feed_id = f.id
    JOIN feed_follows AS ff ON ff.feed_id = p.feed_id
WHERE
    ff.user_id = $1
ORDER BY
    p.published_at DESC NULLS LAST,
    p.created_at DESC,
    p.id DESC
LIMIT $2;

