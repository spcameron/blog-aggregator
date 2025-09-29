-- name: CreateFeedFollow :one
WITH inserted AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
        VALUES ($1, $2, $3, $4, $5)
    RETURNING
        *)
    SELECT
        i.*,
        u.name AS user_name,
        f.name AS feed_name
    FROM
        inserted AS i
            INNER JOIN users AS u ON i.user_id = u.id
            INNER JOIN feeds AS f ON i.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    u.name AS user_name,
    f.name AS feed_name,
    f.url AS feed_url
FROM
    feed_follows AS ff
    JOIN users AS u ON ff.user_id = u.id
    JOIN feeds AS f ON ff.feed_id = f.id
WHERE
    ff.user_id = $1
ORDER BY
    f.name,
    f.id;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE feed_id = $1
    AND user_id = $2;

