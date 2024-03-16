package goob

import (
	"bytes"
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type Room struct {
	Id          uint64
	Messages    []Message
	Users       []*User
	RecvChannel chan Message
	expiryTimer *time.Timer
}

var roomList = make(map[uint64]*Room)

var adjectiveList = [...]string{
	"Astute",
	"Brazen",
	"Charming",
	"Devious",
	"Elderly",
	"Functional",
	"Glorious",
	"Hilarious",
	"Indifferent",
	"Jacked",
	"Knowledgeable",
	"Lazy",
	"Manic",
	"Narcisstic",
	"Obnoxious",
	"Polite",
	"Quirky",
	"Rambunctious",
	"Savant",
	"Tropical",
	"Useless",
	"Valliant",
	"Woke",
	"Xenophobic",
	"Young",
	"Zealous",
}

var animalList = [...]string{
	"Antelope",
	"Beetle",
	"Camel",
	"Dolphin",
	"Elephant",
	"Frog",
	"Giraffe",
	"Hippo",
	"Iguana",
	"Jellyfish",
	"Koala",
	"Lobster",
	"Meerkat",
	"Narwal",
	"Orangutan",
	"Pinguin",
	"Quetzal",
	"Rhino",
	"Shark",
	"Turtle",
	"Unicorn",
	"Viper",
	"Whale",
	"Yak",
	"Zebra",
}

var roomExpiryTimerDuration = RoomExpiryTimeTimerDuration() * time.Second

func GetRoom(id uint64) *Room {
	room, ok := roomList[id]

	if !ok {
		return nil
	}

	return room
}

func NewRoom() *Room {
	id := rand.Uint64() % 1_000_000

	for GetRoom(id) != nil {
		id = rand.Uint64() % 1_000_000
	}

	newRoom := &Room{
		Id:          id,
		Messages:    make([]Message, 0, 10),
		Users:       make([]*User, 0, 10),
		RecvChannel: make(chan Message),
		expiryTimer: nil,
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

		regularMessageBuffer := bytes.Buffer{}
		err := MessageComponent(msg, false).Render(context.Background(), &regularMessageBuffer)

		if err != nil {
			log.Printf("Room %d | Shat my pants when trying to render regular message %s", room.Id, msg.Content)
			continue
		}

		ownMessageBuffer := bytes.Buffer{}
		err = MessageComponent(msg, true).Render(context.Background(), &ownMessageBuffer)

		if err != nil {
			log.Printf("Room %d | Shat my pants when trying to render own message: %s", room.Id, msg.Content)
			continue
		}

		for _, user := range room.Users {
			if user == msg.User {
				user.RoomRecvChannel <- ownMessageBuffer.Bytes()
			} else {
				user.RoomRecvChannel <- regularMessageBuffer.Bytes()
			}
		}
	}
}

func (room *Room) AcceptConn(w http.ResponseWriter, r *http.Request) {

	// Stop the autoCloseTimer if it's ticking
	if room.expiryTimer != nil {
		// If we failed to stop it, the room is closed, reject connection
		if !room.expiryTimer.Stop() {
			return
		}
		log.Printf("Room %d | Deletion cancelled", room.Id)
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user := NewUser(r.Context(), room.generateUniqueName(), room, conn)

	room.Users = append(room.Users, user)

	log.Printf("Room %d | New user: [%v] User count: %d ", room.Id, user, len(room.Users))
	go room.sendUserJoinNotification(user)
	user.Runner()
	room.RemoveUser(user)
}

func (room *Room) RemoveUser(user *User) {

	for i, u := range room.Users {
		if u == user {
			room.Users[i] = room.Users[len(room.Users)-1]
			room.Users = room.Users[:len(room.Users)-1]
		}
	}

	room.sendUserDisconnectNotification(user)
	log.Printf("Room %d | Disconnected user: [%v] User count: %d ", room.Id, user, len(room.Users))

	if len(room.Users) == 0 {
		log.Printf("Room %d | Room empty, expiring in %v", room.Id, roomExpiryTimerDuration)
		room.expiryTimer = time.AfterFunc(roomExpiryTimerDuration, func() {
			log.Printf("Room %d | Deleted", room.Id)
			roomList[room.Id] = nil
		})
	}
}

func (room *Room) sendUserJoinNotification(user *User) {

	userJoinNotificationBuffer := bytes.Buffer{}
	err := UserJoinComponent(user).Render(context.Background(), &userJoinNotificationBuffer)

	if err != nil {
		log.Printf("Room %d | Failed to render user join notification User [%v]", room.Id, user)
		return
	}

	for _, user := range room.Users {
		user.RoomRecvChannel <- userJoinNotificationBuffer.Bytes()
	}
}

func (room *Room) sendUserDisconnectNotification(user *User) {

	userDisconnectNotificationBuffer := bytes.Buffer{}
	err := UserDisconnectComponent(user).Render(context.Background(), &userDisconnectNotificationBuffer)

	if err != nil {
		log.Printf("Failed to render user disconnect notification")
		return
	}

	for _, user := range room.Users {
		user.RoomRecvChannel <- userDisconnectNotificationBuffer.Bytes()
	}
}

func (room *Room) generateUniqueName() UserName {
	var i int
	for {
		i = rand.Int() % len(adjectiveList)
		alreadyExists := false
		for _, user := range room.Users {
			if user.Name[0][0] == adjectiveList[i][0] {
				alreadyExists = true
			}
		}
		if !alreadyExists {
			break
		}
	}

	var j int
	for {
		j = rand.Int() % len(animalList)
		alreadyExists := false
		for _, user := range room.Users {
			if user.Name[1][0] == animalList[j][0] {
				alreadyExists = true
			}
		}
		if !alreadyExists {
			break
		}
	}
	return UserName{adjectiveList[i], animalList[j]}
}
