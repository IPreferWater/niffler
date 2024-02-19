package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	previousSignal = -999
)

type progressMsg struct {
	rssi     int
	rssiDiff int
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

// Exercise: Create a program with two goroutines.
// One goroutine should send "ping" to a channel and the other "pong" in an infinite loop.
// Use a select statement to receive and print messages from the channel.
func main() {
	c := make(chan string)


	go func() {
		for {
			time.Sleep(time.Second*5)
			c <- "ping"
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second*2)
			c <- "pong"
			close(c)
		}
	}()

	for {
		select {
		case x := <-c:
			fmt.Println(x)
		default:
			// Do nothing, just continue the loop
		}
	}
}
func mainA() {

	ctx, cancel := context.WithCancel(context.Background())
	nc, errConnect := nats.Connect("nats://localhost:4222")
	if errConnect != nil {
		panic(errConnect)
	}
	go func() {
		// listen for interrupts to exit gracefully
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
		<-sigChannel
		close(sigChannel)
		cancel()
	}()

	go consumer(ctx, nc)
	startGraphical()
	<-ctx.Done()

	log.Println("server shutdown completed")
	log.Println("exiting gracefully")

}

func consumer(ctx context.Context, nc *nats.Conn) {
	messages := make(chan *nats.Msg, 1000)
	subject := "niffler.rfid"

	// we're subscribing to the subject
	// and assigning our channel as reference to receive messages there
	subscription, err := nc.ChanSubscribe(subject, messages)
	if err != nil {
		log.Fatal("Failed to subscribe to subject:", err)
	}

	defer func() {
		subscription.Unsubscribe()
		close(messages)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from consumer")
			return
		case msg := <-messages:

			split := strings.Split(string(msg.Data), ",")
			_, rssiStr := split[0], split[1]

			rssi, err := strconv.Atoi(rssiStr)
			if err != nil {
				log.Fatal(err)
			}

			if previousSignal == -999 {
				//do nothing as it's the first one
				//log.Printf("first time we got rssi signal, new rssi is set to %d",rssi)
				previousSignal = rssi
				continue
			}
			rssiDiff := rssi - previousSignal

			p.Send(progressMsg{
				rssi:     rssi,
				rssiDiff: rssiDiff,
			})
			previousSignal = rssi
		}
	}
}

func getMessageForRssiSignal(diff int) string {
	if diff > 0 {
		return "you are getting closer"
	}
	if diff < 0 {
		return "you are getting further"
	}

	return "nothing changed"
}
