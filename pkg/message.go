package goob

import "time"

type Message struct {
	Content string
	Time    time.Time
}

func NewMessage(content string) Message {
	return Message{
		Content: content,
		Time:    time.Now(),
	}
}
