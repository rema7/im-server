package hub

type Message struct {
	Type      string `json:"type,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

type Packet struct {
	Token   string  `json:"token"`
	Message Message `json:"message,omitempty"`
}
