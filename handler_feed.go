// handler_feed.go

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spcameron/blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("%d args passed, but `addfeed` expects two arguments, a name and a URL", len(cmd.Args))
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	now := time.Now().UTC()

	feed, err := s.db.CreateFeed(
		ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      feedName,
			Url:       feedURL,
			UserID:    user.ID,
		},
	)

	fmt.Printf("feed created: %s (%s)", feed.Name, feed.ID)
	return nil

}
