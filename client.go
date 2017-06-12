package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

type Client struct {
	id     string
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.unregister <- c
			c.conn.Close()
			break
		}
		jsonMessage, _ := json.Marshal(Message{Sender: c.id, Content: string(message)})
		c.hub.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
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