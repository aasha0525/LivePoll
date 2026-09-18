package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"livepoll/db"
)

var upgrader = websocket.Upgrader{
	// Any origin is allowed here for simplicity; tighten this for production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// PollWebSocket upgrades the connection and streams live vote updates for
// one poll to the client, by subscribing to that poll's Redis Pub/Sub channel.
func PollWebSocket(c *gin.Context) {
	pollID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	sub := db.Rdb.Subscribe(db.Ctx, db.PollChannel(pollID))
	defer sub.Close()

	ch := sub.Channel()

	// Detect client disconnects so we can stop the goroutine cleanly.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}
