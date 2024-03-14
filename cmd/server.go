package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"hugobde.dev/goob/pkg"
)

func main() {
	fmt.Println("Starting...")
	fs := http.FileServer(http.Dir("./public"))

	// Static files serving
	http.Handle("/public/", http.StripPrefix("/public", fs))

	http.Handle("/", templ.Handler(goob.HomeTemplate()))

	http.HandleFunc("/newroom", newRoomHandler)

	http.HandleFunc("/room", roomHandler)

	http.HandleFunc("/chat/{roomId}", chatHandler)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "pooper", http.StatusBadRequest)
		return
	}

	roomId, err := strconv.ParseUint(r.FormValue("id"), 10, 64)

	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	room := goob.GetRoom(roomId)

	if room == nil {
		http.NotFound(w, r)
		return
	}

	goob.RoomTemplate(room).Render(r.Context(), w)
}

func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := goob.NewRoom()
	roomURL := fmt.Sprintf("/room?id=%d", room.Id)
	http.Redirect(w, r, roomURL, http.StatusFound)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	roomIdStr := r.PathValue("roomId")

	if roomIdStr == "" {
		fmt.Println("Missing room ID in /chat somehow")
		return
	}

	roomId, err := strconv.ParseUint(roomIdStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	room := goob.GetRoom(roomId)
	if room == nil {
		log.Printf("Invalid room ID in /chat somehow %s", roomIdStr)
		return
	}

	room.AcceptConn(w, r)
}
