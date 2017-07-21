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

func (Client) decodeMessage(data []byte) (RequestMessage, error) {
	var p RequestMessage
	err := json.Unmarshal(data, &p)

	if err != nil {
		return RequestMessage{}, err
	}

	return p, nil
}

func (c Client) encodeMessage(messageType string, payload interface{}) ([]byte, error) {
	result, err := json.Marshal(ResponseMessage{
		Type:    messageType,
		Payload: payload,
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

		message, err := c.decodeMessage(data)
		if err != nil {
			println(err)
		}

		redis := cache.GetCache()
		senderID, err := redis.GetUserId(message.Token)
		var jsonMessage []byte
		if err != nil {
			jsonMessage, _ = c.encodeMessage("ERROR_MESSAGE", ErrorMessage{
				Content: "No id in cache",
			})
			c.hub.unregister <- c
		} else {
			c.id = senderID
			if message.Type == "CHAT_MESSAGE" {
				var chatMessage ChatMessage
				err := json.Unmarshal([]byte(message.Payload), &chatMessage)
				if err != nil {
					jsonMessage, _ = c.encodeMessage("ERROR_MESSAGE", ErrorMessage{
						Content: "Bad CHAT_MESSAGE body",
					})
				} else {
					jsonMessage, _ = c.encodeMessage(message.Type, ChatMessage{
						Sender:  senderID,
						ChatID:  chatMessage.ChatID,
						Content: chatMessage.Content,
					})
				}
			}
		}
		c.hub.broadcast <- Result{c.id, jsonMessage}
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
