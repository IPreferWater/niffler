package main

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

var natsServer *nats.Conn

func connectNats() {
	nc, errConnect := nats.Connect("nats://localhost:4222")
	if errConnect != nil {
		panic(errConnect)
	}
	natsServer = nc
}


func publish(rssi int) error {
	message := fmt.Sprintf("%s,%d", "randomEpc", rssi)
	return natsServer.Publish("niffler.rfid", []byte(message))
}
