package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi/api"
	"github.com/fzdwx/jacksapi/cb"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	client = api.NewClient(os.Getenv("EMM_API_KEY"))

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
				[]api.ChatMessage{
					{Role: "user", Content: content},
				}).
				DoWithCallback(cb.Output)
		}}
)

func main() {
	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
