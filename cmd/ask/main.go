package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi"
	"github.com/spf13/cobra"
	"os"
)

var (
	client = ai.NewClient(os.Getenv("EMM_API_KEY"))

	root = cobra.Command{
		Use:   "ask <question>",
		Short: "Ask a question",
	}
)

func main() {
	if len(os.Args) == 0 || os.Args[1] != "server" {
		ask()
		return
	}

	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
