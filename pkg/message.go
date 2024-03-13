package goob

import "time"

type Message struct {
	User    *User
	Content string
	Time    time.Time
}

func NewMessage(user *User, content string) Message {
	return Message{
		User:    user,
		Content: content,
		Time:    time.Now(),
	}
}
