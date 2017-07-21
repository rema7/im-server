package hub

import "encoding/json"

type ErrorMessage struct {
	Content string `json:"content,omitempty"`
}

type ChatMessage struct {
	ChatID  int    `json:"chat_id"`
	Sender  int    `json:"sender_id,omitempty"`
	Content string `json:"content,omitempty"`
}

type RequestMessage struct {
	Token   string          `json:"token,omitempty"`
	Type    string          `json:"type,omitempty"`
	Payload json.RawMessage `json:"payload"`
}

type ResponseMessage struct {
	Type    string      `json:"type,omitempty"`
	Payload interface{} `json:"payload"`
}
