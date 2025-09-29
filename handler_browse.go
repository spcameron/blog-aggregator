// handler_browse.go

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spcameron/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("%d args passed, but `browse` expects at most one optional argument, a limit parameter", len(cmd.Args))
	}

	limit := int32(2)
	if len(cmd.Args) == 1 {
		n, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil || n <= 0 {
			return fmt.Errorf("invalid limit %q (must be a positive integer)", cmd.Args[0])
		}
		limit = int32(n)
	}

	ctx := context.Background()

	posts, err := s.db.GetPostsForUser(
		ctx,
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return fmt.Errorf("get posts for user: %w", err)
	}

	printPosts(posts)
	return nil
}

func printPosts(posts []database.GetPostsForUserRow) {
	if len(posts) == 0 {
		fmt.Println("No posts found. Try adding more feeds and scraping again.")
		return
	}

	for i, p := range posts {
		desc := "(no description)"
		if p.Description.Valid {
			desc = p.Description.String
		}

		fmt.Printf("[Post published:  %s]\n", p.PublishedAt.Format(time.RFC3339))
		fmt.Printf("* Title:          %s\n", p.Title)
		fmt.Printf("* Feed Name:      %s\n", p.FeedName)
		fmt.Printf("* URL:            %s\n", p.Url)
		fmt.Printf("* Description:    %s\n", desc)

		if i < len(posts)-1 {
			fmt.Println()
		}
	}
}
