package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marcel-alter/AggreGator/internal/config"
)

func main() {
	fmt.Printf("Reading from filePath ~/.gatorconfig.json\n")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error: Reading failed: %v", err)
		return
	}
	sta := state{cfg: &cfg}
	comms := commands{registeredCommands: make(map[string]func(*state, command) error)}
	comms.register("login", handlerLogin)
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
