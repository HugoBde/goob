package goob

import (
	"context"
	"encoding/json"
	"log"

	"nhooyr.io/websocket"
)

type WSMessage struct {
	Message string `json:"message"`
}

type User struct {
	ctx             context.Context
	conn            *websocket.Conn
	wsSendChannel   chan []byte
	wsRecvChannel   chan []byte
	RoomRecvChannel chan []byte
	roomSendChannel chan Message
}

func NewUser(ctx context.Context, room *Room, conn *websocket.Conn) *User {
	return &User{
		ctx:             ctx,
		conn:            conn,
		wsSendChannel:   make(chan []byte),
		wsRecvChannel:   make(chan []byte),
		RoomRecvChannel: make(chan []byte),
		roomSendChannel: room.RecvChannel,
	}
}

func (user *User) websocketReader() {
	for {
		_, data, err := user.conn.Read(user.ctx)
		if err != nil {
			close(user.wsRecvChannel)
			return
		}
		user.wsRecvChannel <- data
	}
}

func (user *User) websocketWriter() {
	for {
		data, ok := <-user.wsSendChannel
		if !ok {
			return
		}
		err := user.conn.Write(user.ctx, websocket.MessageText, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (user *User) Runner() {
	defer close(user.wsSendChannel)
	defer user.conn.CloseNow()

	go user.websocketReader()
	go user.websocketWriter()

	for {
		select {
		case data, ok := <-user.wsRecvChannel:
			if !ok {
				return
			}
			var wsMessage WSMessage
			err := json.Unmarshal(data, &wsMessage)
			if err != nil {
				log.Println(err)
				return
			}
			msg := NewMessage(user, wsMessage.Message)
			user.roomSendChannel <- msg
		case data, ok := <-user.RoomRecvChannel:
			if !ok {
				return
			}
			user.wsSendChannel <- data
		}
	}
}
