package goob

import (
	"context"
	"log"

	// "math/rand"
	"bytes"
	"net/http"

	"nhooyr.io/websocket"
)

type Room struct {
	Id          uint64
	Messages    []Message
	Users       []*User
	RecvChannel chan Message
}

var roomList = make(map[uint64]*Room)

func GetRoom(id uint64) *Room {
	room, ok := roomList[id]

	if !ok {
		return nil
	}

	return room
}

func NewRoom() *Room {
	// id := rand.Uint64() % 100
	id := uint64(1)

	newRoom := &Room{
		Id:          id,
		Messages:    make([]Message, 0, 10),
		Users:       make([]*User, 0, 10),
		RecvChannel: make(chan Message),
	}

	roomList[id] = newRoom

	go newRoom.Run()

	return newRoom
}

func (room *Room) Run() {
	log.Printf("Room %d running", room.Id)
	for {
		msg := <-room.RecvChannel

		log.Printf("Room %d | #%d | \"%s\"", room.Id, len(room.Messages)+1, msg.Content)

		room.Messages = append(room.Messages, msg)

		buf := bytes.Buffer{}
		err := MessageComponent(msg).Render(context.Background(), &buf)
		if err != nil {
			log.Printf("Room %d | Shat my pants when trying to render %s", room.Id, msg.Content)
			continue
		}

		for _, user := range room.Users {
			user.RoomRecvChannel <- buf.Bytes()
		}
	}
}

func (room *Room) AcceptConn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Room %d | New User [%s]", room.Id, r.Host)

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user := NewUser(r.Context(), room, conn)

	room.Users = append(room.Users, user)
	log.Printf("Room %d | User count: %d ", room.Id, len(room.Users))

	user.Runner()

	room.RemoveUser(user)
}

func (room *Room) RemoveUser(user *User) {

	log.Printf("Room %d | Removing User", room.Id)

	for i, u := range room.Users {
		if u == user {
			room.Users[i] = room.Users[len(room.Users)-1]
			room.Users = room.Users[:len(room.Users)-1]
		}
	}
	log.Printf("Room %d | User count: %d ", room.Id, len(room.Users))
}
