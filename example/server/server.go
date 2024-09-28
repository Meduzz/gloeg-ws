package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/Meduzz/gloegg-ws/api"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func main() {
	srv := gin.Default()

	feedback := make(chan *api.WsEvent, 10)
	pusher := make(chan *api.WsEvent, 10)

	go func() {
		for event := range feedback {
			if event.Kind == api.EventKindLog {
				fmt.Printf("Log: %s\n", string(event.Log))
			} else if event.Kind != api.EventKindFlagRequest {
				fmt.Printf("Flag: %s %s %s %v\n", event.Kind, event.Flag.Name, event.Flag.Kind, event.Flag.Value)
			} else {
				fmt.Printf("Flags where requested\n")
			}
		}
	}()

	srv.GET("/ws", gin.WrapH(websocket.Handler(handler(feedback, pusher))))
	srv.POST("/ws", func(ctx *gin.Context) {
		event := &api.WsEvent{}
		err := ctx.BindJSON(event)

		if err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		pusher <- event
		ctx.Status(201)
	})

	srv.Run(":8080")
}

func handler(feedback chan *api.WsEvent, push chan *api.WsEvent) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		go func() {
			for it := range push {
				err := websocket.JSON.Send(ws, it)

				if err != nil && errors.Is(err, io.EOF) {
					break
				}
			}
		}()

		for {
			event := &api.WsEvent{}
			err := websocket.JSON.Receive(ws, event)

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			} else {
				feedback <- event
			}
		}
	}
}
