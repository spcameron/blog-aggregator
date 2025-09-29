// handler_agg.go

package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	if err := db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: time.Now().UTC(),
		},
	); err != nil {
		return fmt.Errorf("mark feed fetched: %w", err)
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("rss fetch feed: %w", err)
	}

	if len(feedData.Channel.Item) == 0 {
		log.Printf("Feed %s has no new posts since last fetch", feed.Name)
		return nil
	}

	printFeedTitles(feedData)
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	return nil
}

func printFeedTitles(feed *rss.RSSFeed) {
	items := feed.Channel.Item

	if len(items) == 0 {
		fmt.Println("There are no items to print in this feed.")
	}

	for _, item := range items {
		fmt.Printf("New post: %s\n", item.Title)
	}
}
