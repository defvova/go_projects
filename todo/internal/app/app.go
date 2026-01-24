package app

import (
	"fmt"
	"os"
	"todo/internal/cli"
)

func Run(args []string) {
	if len(args) < 2 {
		usage()
	}

	cmd := args[1]
	switch cmd {
	case "ls":
		cli.Ls(args[2:])
	case "add":
		cli.Add(args[2:])
	case "rm":
		cli.Rm(args[2:])
	case "complete":
		cli.Complete(args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Println("usage: tasks <ls<-a>|add<desc>|rm<id>|complete<id>>")
	os.Exit(1)
}
