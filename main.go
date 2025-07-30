package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	config "github.com/saifullah605/Gator/internal/config"
	"github.com/saifullah605/Gator/internal/database"
)

func main() {
	currConfig, err := config.Read()
	if err != nil {
		fmt.Println("Cannot read config:", err)
	}

	db, err := sql.Open("postgres", currConfig.DBURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	states := &state{dbQueries, &currConfig}
	commands := &commands{make(map[string]func(*state, command) error)}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)

	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Println("invalid command, not enough arguments")
		os.Exit(1)
	}

	command := command{
		name:      arguments[1],
		arguments: arguments[2:],
	}

	if err := commands.run(states, command); err != nil {
		fmt.Println("erorr:", err)
		os.Exit(1)
	}

	os.Exit(0)

}
