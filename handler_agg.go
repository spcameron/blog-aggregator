// handler_agg.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/spcameron/blog-aggregator/internal/database"
	"github.com/spcameron/blog-aggregator/rss"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("%d args passed, but `agg` expects one argument, the time between requests", len(cmd.Args))
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s ...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("scrape feeds: %w", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("get next feed: %w", err)
	}
	if err := scrapeFeed(s.db, feed); err != nil {
		return fmt.Errorf("scrape feed: %w", err)
	}

	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) error {
	ctx := context.Background()
	now := time.Now().UTC()

	feedData, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("rss fetch feed: %w", err)
	}

	if err := db.MarkFeedFetched(
		ctx,
		database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: now,
		},
	); err != nil {
		return fmt.Errorf("mark feed fetched: %w", err)
	}

	skipped := 0
	for _, item := range feedData.Channel.Item {

		if _, err := db.CreatePost(
			ctx,
			database.CreatePostParams{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     item.Title,
				Url:       item.Link,
				Description: sql.NullString{
					String: item.Description,
					Valid:  item.Description != "",
				},
				PublishedAt: parsePubDate(item.PubDate),
				FeedID:      feed.ID,
			},
		); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				skipped += 1
			} else {
				log.Println(err)
			}
		}
	}

	log.Printf("%v posts skipped because of duplicate URLs", skipped)
	return nil
}

func parsePubDate(s string) time.Time {
	if s == "" {
		return time.Now().UTC()
	}

	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}

	return time.Now().UTC()
}
