package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"hugobde.dev/goob/pkg"
)

func main() {
	var helpFlag = flag.Bool("help", false, "print help message")
	var certFileFlag = flag.String("cert", "", "TLS domain cert file")
	var keyFileFlag = flag.String("key", "", "TLS private key file")
	var portFlag = flag.Uint("port", 42069, "port to listen on")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if (*certFileFlag != "") != (*keyFileFlag != "") {
		flag.PrintDefaults()
		return
	}

	fs := http.FileServer(http.Dir("./public"))

	// Static files serving
	http.Handle("/public/", http.StripPrefix("/public", fs))

	http.Handle("/", templ.Handler(goob.HomeTemplate()))

	http.HandleFunc("/newroom", newRoomHandler)

	http.HandleFunc("/room", roomHandler)

	http.HandleFunc("/chat/{roomId}", chatHandler)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	/* If a domainCertFlag flag is configured then enable TLS */
	if *certFileFlag != "" {
		log.Printf("Listening on %d with TLS enabled ...", *portFlag)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *portFlag),
			*certFileFlag,
			*keyFileFlag,
			nil))
	} else {
		log.Printf("Listening on %d with TLS disabled ...", *portFlag)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil))
	}
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

	if len(room.Users) >= 10 {
		http.Error(w, "Room full", http.StatusForbidden)
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

	if len(room.Users) >= 10 {
		http.Error(w, "Room full", http.StatusForbidden)
		return
	}

	room.AcceptConn(w, r)
}
