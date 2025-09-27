// handler_reset.go

package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("%d args passed, but `reset` expects zero arguments", len(cmd.Args))
	}

	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("reset users: %w", err)
	}

	log.Printf("users table successfully truncated")
	return nil

}
