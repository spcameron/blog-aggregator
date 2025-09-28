// main.go

package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spcameron/blog-aggregator/internal/config"
	"github.com/spcameron/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Read() error: %v", err)
	}

	dbURL := cfg.URL

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection, %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerListFeedFollows)

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
