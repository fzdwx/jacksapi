package main

import (
	"encoding/json"
	"fmt"
	"github.com/fzdwx/jacksapi/api"
	"github.com/fzdwx/jacksapi/cb"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server [port]",
		Short: "provide chatGpt api",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please input your port")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			g := gin.Default()
			g.POST("/", handle)

			g.POST("/v1/chat/completions", handleMockOpenAi)

			g.Run(":" + args[0])
		},
	}
)

func init() {
	root.AddCommand(serverCmd)
}

func handle(c *gin.Context) {
	var body api.Body
	err := c.ShouldBind(&body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(c.Writer, "Error: %s", err.Error())
		return
	}

	client.ChatStream(body.Messages).
		Temperature(body.Temperature).
		PresencePenalty(body.PresencePenalty).
		DoWithCallback(func(resp *http.Response, err error) {
			io.Copy(c.Writer, resp.Body)
		})

}

func handleMockOpenAi(c *gin.Context) {
	var body api.Body
	err := c.ShouldBind(&body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(c.Writer, "Error: %s", err.Error())
		return
	}

	id := "jacks-" + time.Now().Format("20060102150405")
	created := time.Now().UnixMilli()

	c.Stream(func(w io.Writer) bool {
		client.ChatStream(body.Messages).
			Temperature(body.Temperature).
			PresencePenalty(body.PresencePenalty).
			DoWithCallback(cb.With(func(r rune, done bool, err error) {
				if done || err != nil {
					c.SSEvent("", " [DONE]")
				} else {
					c.SSEvent("", " "+buildMessage(r, id, created, body.Model))
					c.Writer.Flush()
				}
			}))
		return false
	})
}

func buildMessage(r rune, id string, created int64, model string) string {
	return marshal(map[string]interface{}{
		"id":      id,
		"object":  "chat.completion.chunk",
		"created": created,
		"model":   model,
		"choices": []map[string]interface{}{
			{
				"delta": map[string]interface{}{
					"content": string(r),
				},
				"index":         0,
				"finish_reason": nil,
			},
		},
	})
}

func marshal(m map[string]interface{}) string {
	s, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(s)
}
