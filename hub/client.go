package hub

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Print("disconnected")
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.unregister <- c
			c.conn.Close()
			break
		}
		jsonMessage, _ := json.Marshal(Message{Type: "CHAT_MESSAGE", Sender: c.id, Content: string(message)})
		c.hub.broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
