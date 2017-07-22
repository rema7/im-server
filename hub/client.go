package hub

import (
	"encoding/json"
	"im-server/db"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int64
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

		var jsonMessage []byte
		var memberIds = []int{}
		if err != nil {
			jsonMessage, _ = c.encodeMessage("ERROR_MESSAGE", ErrorMessage{
				Content: "No id in cache",
			})
			c.hub.unregister <- c
		} else {
			if message.Type == "CHAT_MESSAGE" {
				var chatMessage ChatMessage
				err := json.Unmarshal([]byte(message.Payload), &chatMessage)
				if err != nil {
					jsonMessage, _ = c.encodeMessage("ERROR_MESSAGE", ErrorMessage{
						Content: "Bad CHAT_MESSAGE body",
					})
				} else {
					session := db.GetDbSession()
					var chatID int
					chatMember := db.ChatMember{}
					err := session.Model(&chatMember).Column("chat_id").Where("chat_id = ?", chatMessage.ChatID).Where("user_id = ?", c.id).Select(&chatID)
					if err != nil {
						log.Println(err)
					}
					err = session.Model(&chatMember).Column("user_id").Where("chat_id = ?", chatMessage.ChatID).Where("user_id != ?", c.id).Select(&memberIds)
					if err != nil {
						log.Println(err)
						return
					}
					jsonMessage, _ = c.encodeMessage(message.Type, ChatMessage{
						Sender:  c.id,
						ChatID:  chatMessage.ChatID,
						Content: chatMessage.Content,
					})
				}
			}
		}
		c.hub.broadcast <- Result{memberIds, jsonMessage}
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
