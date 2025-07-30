package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	config "github.com/saifullah605/Gator/internal/config"
	"github.com/saifullah605/Gator/internal/database"
)

type state struct {
	db     *database.Queries
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

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])

	if err == sql.ErrNoRows {
		return fmt.Errorf("user does not exist")
	} else if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	err = s.config.SetUser(cmd.arguments[0])

	if err != nil {
		return fmt.Errorf("cannot set username: %v", err)
	}

	fmt.Println("User has been set:", cmd.arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the command neeeds a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		return fmt.Errorf("cannot create user, name already used")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error: %v", err)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})

	if err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	fmt.Println("User", user.Name, "created successfully")

	return s.config.SetUser(user.Name)

}

func handlerReset(s *state, cmd command) error {
	return nil
}
