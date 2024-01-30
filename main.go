package main

import (
	"Go-Tutorials/Core-lang/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Core programming language!\n",
		user.Username)
	fmt.Printf("You can type now\n")
	repl.Start(os.Stdin, os.Stdout)
}
