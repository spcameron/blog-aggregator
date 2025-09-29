# Blog Aggregator (CLI)

This is the sixth project on Boot.dev. A fast, minimal RSS/Atom feed aggregator written in Go. It lets you follow feeds, scrape new posts into PostgreSQL, and browse them from your terminal.

---

## Prerequisites

* **Go** ≥ 1.22 (tested with 1.25.x)
* **PostgreSQL** ≥ 15

---

## Install

The module path is:

```
module github.com/spcameron/blog-aggregator
```

Because the `main` package lives at the repo root, the installed binary name is **`blog-aggregator`**. Install it with:

```bash
go install github.com/spcameron/blog-aggregator@latest
```

After installation, run the CLI with:

```bash
blog-aggregator <command> [args]
```

---

## Configuration

Create a JSON config file at **`~/.gatorconfig.json`**:

```json
{
  "db_url": "postgres://USERNAME:PASSWORD@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your-username"
}
```

* `db_url` — PostgreSQL connection string
* `current_user_name` — default user context for CLI operations

---

## Database Setup

Goose migration tooling is provided:

Apply the schema in order (from the repo root):

```bash
psql "$DB_URL" -f sql/schema/001_users.sql
psql "$DB_URL" -f sql/schema/002_feeds.sql
psql "$DB_URL" -f sql/schema/003_feed_follows.sql
psql "$DB_URL" -f sql/schema/004_feeds.sql
psql "$DB_URL" -f sql/schema/005_posts.sql
```

---

## Usage

Implemented commands (source in `commands.go` and `handler_*.go`):

* **User context**

  * `register <name>` — add a new user 
  * `login <name>` - set the current user in config
  * `users` - list all users in the database

* **Feeds**

  * `addfeed <url> <name>` — add a feed to the catalog
  * `feeds` — list all known feeds
  * `follow <url>` — follow a feed for the current user
  * `following` — list feeds followed by the current user
  * `unfollow <url>` — unfollow a feed

* **Aggregation & browsing**

  * `agg <time_between_fetches>` — fetch new posts across followed feeds at a given interval (e.g. `1m`, `10s`)
  * `browse <limit>` — print stored posts (most recent first)

* **Maintenance (dev only)**

  * `reset` — truncate/clear development data

---

## Quick Start

```bash
# 0) Ensure Postgres is running and the schema is applied (see above)

# 1) Point the CLI at your DB and user
printf '{
  "db_url": "postgres://spc:@localhost:5432/gator?sslmode=disable",
  "current_user_name": "sean"
}\n' > ~/.gatorconfig.json

# 2) Add + follow a feed
blog-aggregator addfeed <feed-name> https://example.com/rss 
blog-aggregator follow https://example.com/rss

# 3) Aggregate new posts every 1 minute
blog-aggregator agg 1m

# 4) Browse stored posts (limit 25)
blog-aggregator browse 25
```

Typical `browse` output shows each post’s published date, title, feed name, URL, and a description.

---

## Project Structure

```
.
├── commands.go
├── handler_agg.go
├── handler_browse.go
├── handler_feed.go
├── handler_follow.go
├── handler_reset.go
├── handler_user.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   └── database/
│       ├── db.go
│       ├── feed_follows.sql.go
│       ├── feeds.sql.go
│       ├── models.go
│       ├── posts.sql.go
│       └── users.sql.go
├── rss/
│   └── rss.go
├── sql/
│   ├── queries/
│   │   ├── feed_follows.sql
│   │   ├── feeds.sql
│   │   ├── posts.sql
│   │   └── users.sql
│   └── schema/
│       ├── 001_users.sql
│       ├── 002_feeds.sql
│       ├── 003_feed_follows.sql
│       ├── 004_feeds.sql
│       └── 005_posts.sql
├── main.go
├── go.mod
├── go.sum
└── README.md
```

---

## License

MIT License. See `LICENSE` for details.

---

## Repository

[https://github.com/spcameron/blog-aggregator](https://github.com/spcameron/blog-aggregator)
