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
			g.POST("/", func(context *gin.Context) {
				var body api.Body
				err := context.ShouldBind(&body)
				if err != nil {
					context.Writer.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(context.Writer, "Error: %s", err.Error())
					return
				}

				var (
					rchan    = make(chan rune)
					doneChan = make(chan bool)
				)

				go func() {
					client.ChatStream(body.Messages).
						Temperature(body.Temperature).
						PresencePenalty(body.PresencePenalty).
						DoWithCallback(cb.With(func(r rune, done bool, err error) {
							if done || err != nil {
								doneChan <- true
								return
							}
							rchan <- r
						}))
				}()

				context.Stream(func(w io.Writer) bool {
					select {
					case r := <-rchan:
						fmt.Println("rune:", string(r))
						context.SSEvent("message", string(r))
						return true
					case <-doneChan:
						return false
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
