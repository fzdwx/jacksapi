package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi/api"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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
			g.POST("/", func(c *gin.Context) {
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

			})

			g.Run(":" + args[0])
		},
	}
)

func init() {
	root.AddCommand(serverCmd)
}
