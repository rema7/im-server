package hub

import (
	"encoding/json"
	"im-server/cache"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (Client) decodeMessage(data []byte) (Packet, error) {
	var p Packet
	err := json.Unmarshal(data, &p)

	if err != nil {
		return Packet{}, err
	}

	return p, nil
}

func (c Client) encodeMessage(messageType string, message Message) ([]byte, error) {
	result, err := json.Marshal(Message{
		Type:    messageType,
		ChatID:  message.ChatID,
		Sender:  c.id,
		Content: message.Content,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Print("disconnected")
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.unregister <- c
			c.conn.Close()
			break
		}

		packet, err := c.decodeMessage(data)
		if err != nil {
			println(err)
		}

		redis := cache.Cache{}
		redis.Init()
		senderID, err := redis.GetUserId(packet.Token)
		var jsonMessage []byte

		if err != nil {
			log.Println(err)
			jsonMessage, _ = c.encodeMessage("ERROR_MESSAGE", packet.Message)
		} else {
			c.id = senderID
			jsonMessage, _ = c.encodeMessage("CHAT_MESSAGE", packet.Message)
		}

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
