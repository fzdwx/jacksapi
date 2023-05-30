package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	client = ai.NewClient(os.Getenv("EMM_API_KEY"))

	root = cobra.Command{
		Use:   "ask <question>",
		Short: "Ask a question",
	}
)

func main() {
	if os.Args[1] != "server" {
		ask()
		return
	}

	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ask() {
	if len(os.Args) == 1 {
		fmt.Println("Please input your ask")
		os.Exit(1)
	}

	var (
		content = strings.Join(os.Args[1:], " ")
	)
	client.ChatStream(
		[]ai.ChatMessage{
			{Role: "user", Content: content},
		}).
		DoWithCallback(ai.Output)
}
