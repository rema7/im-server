package hub

type Message struct {
	ChatID  int    `json:"chat_id"`
	Type    string `json:"type,omitempty"`
	Sender  int    `json:"sender_id,omitempty"`
	Content string `json:"content,omitempty"`
}

type Packet struct {
	Token   string  `json:"token"`
	Message Message `json:"message,omitempty"`
}
