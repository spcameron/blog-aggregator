// handler_agg.go

package main

import (
	"context"
	"fmt"

	"github.com/spcameron/blog-aggregator/rss"
)

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
