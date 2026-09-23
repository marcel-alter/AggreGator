package main

import (
	"fmt"

	"github.com/marcel-alter/AggreGator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

//↓↓↓↓↓↓↓↓↓↓↓↓↓ commands methods

func (c *commands) run(s *state, cmd command) error {
	_, exists := c.registeredCommands[cmd.name]
	if exists != false {
		err := c.registeredCommands[cmd.name](s, cmd)
		return err
	} else {
		fmt.Println("commands run err not nil triggered!")
		return fmt.Errorf("Error: Unknown command: %s", cmd.name)
	}
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
