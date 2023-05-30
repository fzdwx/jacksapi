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
		Use:   "ask",
		Short: "Ask a question",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please input your ask")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				content = strings.Join(args, " ")
			)
			client.ChatStream(
				[]ai.ChatMessage{
					{Role: "user", Content: content},
				}).
				DoWithCallback(ai.Output)
		}}
)

func main() {
	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
