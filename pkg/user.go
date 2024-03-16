package goob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"nhooyr.io/websocket"
)

type UserName [2]string

type WSMessage struct {
	Message string `json:"message"`
}

type User struct {
	Name            UserName
	ctx             context.Context
	conn            *websocket.Conn
	room            *Room
	wsSendChannel   chan []byte
	wsRecvChannel   chan []byte
	RoomRecvChannel chan []byte
	roomSendChannel chan Message
}

func NewUser(ctx context.Context, name UserName, room *Room, conn *websocket.Conn) *User {
	return &User{
		Name:            name,
		ctx:             ctx,
		conn:            conn,
		room:            room,
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

	// Closing this connection will end the websocketReader goroutine
	defer user.conn.CloseNow()

	go user.websocketReader()
	go user.websocketWriter()

	// Defer closing the websocketWriter goroutine
	defer close(user.wsSendChannel)

	for {
		select {
		// We receive data from the websocket connection, parse it and forward the message to the room
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

			trimmedMessage := strings.TrimSpace(wsMessage.Message)

			if len(trimmedMessage) != 0 {
				msg := NewMessage(user, trimmedMessage)
				user.roomSendChannel <- msg
			}

			// We receive data from the room, hopefully some html, forward it to the websocket connection
		case data, ok := <-user.RoomRecvChannel:
			if !ok {
				return
			}
			user.wsSendChannel <- data
		}
	}
}

func (user *User) String() string {
	return fmt.Sprintf("%s %s", user.Name[0], user.Name[1])
}
