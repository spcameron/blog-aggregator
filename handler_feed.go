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

	follow, err := s.db.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed follow created successfully:")
	printFeedFollow(follow)
	return nil

}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("%d args passed, but `feeds` expects zero arguments", len(cmd.Args))
	}

	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("%d feeds found:", len(feeds))

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("get user by id: %w", err)
		}

		printFeed(feed, user)
		fmt.Println()
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:         %s\n", feed.ID)
	fmt.Printf("* Created:    %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:    %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:       %s\n", feed.Name)
	fmt.Printf("* URL:        %s\n", feed.Url)
	fmt.Printf("* User:       %s\n", user.Name)
}
