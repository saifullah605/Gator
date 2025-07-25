package main

import (
	"fmt"

	config "github.com/saifullah605/Gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandFunc, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}

	return commandFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the commands needs a username")
	}

	err := s.config.SetUser(cmd.arguments[0])

	if err != nil {
		return fmt.Errorf("cannot set username: %v", err)
	}

	fmt.Println("User has been set:", cmd.arguments[0])

	return nil
}
