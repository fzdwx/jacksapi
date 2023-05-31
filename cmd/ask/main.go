package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi"
	"github.com/spf13/cobra"
	"os"
)

var (
	client = jacksapi.NewClient(os.Getenv(KeyName))

	root = cobra.Command{
		Use:     "ask <question>",
		Short:   "Ask a question",
		Version: Version,
	}
)

func main() {
	if len(os.Args) == 1 || os.Args[1] != "server" {
		ask()
		return
	}

	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
