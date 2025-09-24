// main.go

package main

import (
	"log"
	"os"

	"github.com/spcameron/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Read() error: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}
	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalf("Too few arguments passed, require at least two, but received %d", len(os.Args))
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{cmdName, cmdArgs}

	if err := cmds.run(programState, cmd); err != nil {
		log.Fatalf("%s() returned a non-nil error, %v", cmd.Name, err)
	}
}
