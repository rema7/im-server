package main

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func (self *Message) String() string {
	return self.Email + " says " + self.Message
}
