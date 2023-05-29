package main

import (
	"fmt"
	"github.com/fzdwx/jacksapi/api"
	"github.com/fzdwx/jacksapi/cb"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var (
	proxyCmd = &cobra.Command{
		Use:   "proxy [port]",
		Short: "Proxy ChatGpt api",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please input your port")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			g := gin.Default()
			g.POST("/v1/chat/completions", func(c *gin.Context) {
				var body api.Body
				err := c.ShouldBind(&body)
				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(c.Writer, "Error: %s", err.Error())
					return
				}

				c.Header("Connection", "keep-alive")
				c.Header("Transfer-Encoding", "chunked")
				var (
					rchan    = make(chan rune, 1000)
					doneChan = make(chan bool)
				)

				go client.ChatStream(body.Messages).
					Temperature(body.Temperature).
					PresencePenalty(body.PresencePenalty).
					DoWithCallback(cb.With(func(r rune, done bool, err error) {
						if done || err != nil {
							doneChan <- true
							return
						}
						rchan <- r
					}))

				c.Stream(func(w io.Writer) bool {
					for {
						select {
						case r := <-rchan:
							c.SSEvent("", string(r))
							fmt.Println(string(r))
						case <-doneChan:
							c.SSEvent("", "done")
							return false
						}
					}
				})

			})

			g.Run(":" + args[0])
		},
	}
)

func init() {
	root.AddCommand(proxyCmd)
}
