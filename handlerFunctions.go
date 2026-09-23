package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	//fmt.Printf("LOGIN TRIGGERED, len of args = %d\n", len(cmd.args))
	if len(cmd.args) == 0 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: Login requieres a username argument!")
	}
	/*for _, arg := range cmd.args {
		fmt.Println(arg)
	}*/
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("Error at handlerLogin: %w", err)
	}
	fmt.Printf("User %s was set.", cmd.args[0])
	return nil
}
