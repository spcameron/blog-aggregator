// handler_follow.go

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spcameron/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("%d args passed, but `follow` expects one argument, the URL", len(cmd.Args))
	}

	ctx := context.Background()
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("get feed: %w", err)
	}

	now := time.Now().UTC()

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

	fmt.Println("Feed follow created successfully:")
	printFeedFollow(follow)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("get feed follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}

	return nil
}

func printFeedFollow(follow database.CreateFeedFollowRow) {
	fmt.Printf("* Feed:       %s\n", follow.FeedName)
	fmt.Printf("* User:       %s\n", follow.UserName)
}
