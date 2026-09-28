package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/marcel-alter/AggreGator/internal/config"
	"github.com/marcel-alter/AggreGator/internal/database"
)

func main() {
	//fmt.Printf("Reading from filePath ~/.gatorconfig.json\n")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error: Reading failed: %v", err)
		return
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	defer db.Close()
	sta := state{cfg: &cfg, db: dbQueries}
	comms := commands{registeredCommands: make(map[string]func(*state, command) error)}
	comms.register("login", handlerLogin)
	comms.register("register", handlerRegister)
	comms.register("reset", handlerReset)
	comms.register("users", handlerUsers)
	comms.register("agg", handlerAgg)
	comms.register("addfeed", handlerAddFeed)
	comms.register("feeds", handlerFeeds)
	comms.register("help", handlerHelp)
	comms.register("follow", handlerFollow)
	comms.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		log.Fatal("Error: No Command given! Type 'Help' for List of Commands")
		return
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{name: cmdName, args: cmdArgs}
	if err := comms.run(&sta, cmd); err != nil {
		log.Fatalf("Error when running command: %v ", err)
	}
}
