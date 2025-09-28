-- sql/schema/003_feed_follows.sql
-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;

