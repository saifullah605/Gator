package main

import (
	"fmt"

	config "github.com/saifullah605/Gator/internal/config"
)

func main() {
	currConfig, err := config.Read()
	if err != nil {
		fmt.Println(nil)
	}

	fmt.Println(currConfig)
	if err := config.SetUser("saif"); err != nil {
		fmt.Println(err)
	}

	currConfig, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(currConfig)

}
