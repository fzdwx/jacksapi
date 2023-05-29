package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi/api"
	"github.com/fzdwx/jacksapi/cb"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("Usage: ai <message>")
		return
	}

	var (
		content    = strings.Join(os.Args[1:], " ")
		accessCode = os.Getenv("EMM_API_KEY")
	)
	fmt.Println("Your ask:", content)
	api.NewClient(accessCode).
		ChatStream(
			[]api.ChatMessage{
				{Role: "user", Content: content},
			}).
		DoWithCallback(cb.Output)
}
