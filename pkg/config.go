package goob

import (
	"flag"
	"time"
)

var helpFlag = flag.Bool("help", false, "print help message")
var certFileFlag = flag.String("cert", "", "TLS domain cert file")
var keyFileFlag = flag.String("key", "", "TLS private key file")
var portFlag = flag.Uint("port", 42069, "port to listen on")
var roomExpiryTimeTimerDurationFlag = flag.Uint("room-expiry", 1200, "number of seconds to wait before deleting a room with no users")

func Help() bool {
	if !flag.Parsed() {
		flag.Parse()
	}
	return *helpFlag
}

func CertFile() string {
	if !flag.Parsed() {
		flag.Parse()
	}
	return *certFileFlag
}

func KeyFile() string {
	if !flag.Parsed() {
		flag.Parse()
	}
	return *keyFileFlag
}

func Port() uint {
	if !flag.Parsed() {
		flag.Parse()
	}
	return *portFlag
}

func RoomExpiryTimeTimerDuration() time.Duration {
	if !flag.Parsed() {
		flag.Parse()
	}
	return time.Duration(*roomExpiryTimeTimerDurationFlag)
}
