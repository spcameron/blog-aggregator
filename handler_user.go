// handler_user.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spcameron/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("%d args passed, but `login` expects a single argument, the username", len(cmd.Args))
	}

	name := strings.TrimSpace(cmd.Args[0])
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if _, err := s.db.GetUser(context.Background(), name); err == sql.ErrNoRows {
		return fmt.Errorf("user does not exist: %q", name)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Printf("user has been set successfully to %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("%d args passed, but `register` expects a single argument, the name", len(cmd.Args))
	}

	name := strings.TrimSpace(cmd.Args[0])
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	now := time.Now().UTC()

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      cmd.Args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	fmt.Println()
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, user := range users {
		var currentTag string
		if s.cfg.CurrentUserName == user.Name {
			currentTag = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, currentTag)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf("* ID:         %s\n", user.ID)
	fmt.Printf("* Created:    %v\n", user.CreatedAt)
	fmt.Printf("* Updated:    %v\n", user.UpdatedAt)
	fmt.Printf("* Name:       %s\n", user.Name)
}
