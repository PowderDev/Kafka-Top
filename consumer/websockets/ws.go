package websockets

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 12 * time.Second,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

var writeMutex = &sync.Mutex{}

type WsClient struct {
	hub  *Hub
	conn *websocket.Conn
}

type Message struct {
	Key  string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (c *WsClient) readPump() {
	defer func() {
		c.hub.unregister <- c
	}()
	c.conn.SetReadDeadline(time.Time{})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("received: %s", message)
	}
}

func (c *WsClient) writeMessage(encodedMessage []byte) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	err := c.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return err
	}

	writeMutex.Lock()
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		writeMutex.Unlock()
		return err
	}

	_, err = w.Write(encodedMessage)
	if err != nil {
		writeMutex.Unlock()
		return err
	}

	if err := w.Close(); err != nil {
		writeMutex.Unlock()
		return err
	}

	writeMutex.Unlock()
	return nil
}

func ServeWs(ctx *fasthttp.RequestCtx, hub *Hub) {
	ctx.Request.Header.Set("Connection", "upgrade")
	ctx.Request.Header.Set("Upgrade", "websocket")

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := &WsClient{hub: hub, conn: conn}
		client.hub.register <- client

		client.readPump()
	})

	if err != nil {
		log.Println(err)
	}
}
