// handler_user.go

package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("%d args passed, but `login` expects a single argument, the username", len(cmd.Args))
	}

	user := cmd.Args[0]
	if err := s.cfg.SetUser(user); err != nil {
		return fmt.Errorf("Call to SetUser() failed: %w", err)
	}

	fmt.Printf("User has been set successfully to %s\n", user)
	return nil
}
