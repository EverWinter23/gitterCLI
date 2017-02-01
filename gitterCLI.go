package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sromku/go-gitter"
)

func main() {

	// Setup
	gitter_token := get_env("GITTER_TOKEN")
	client := gitter.New(gitter_token)
	self, err := client.GetUser()
	panic_err(err)

	var room_name string
	fmt.Printf("Enter room name: ")
	fmt.Scanf("%s", &room_name)
	CURRENT_ROOM_ID, err := client.GetRoomId(room_name)
	panic_err(err)

	fmt.Printf("\nSuccessfully joined room " + room_name + "\n")

	// Execution
	receiver := client.Stream(CURRENT_ROOM_ID)
	go client.Listen(receiver)

	sender := make(chan string, 10)
	go get_input(sender)

	// Interaction
	for {
		select {
		case msg := <-receiver.Event:
			switch ev := msg.Data.(type) {
			case *gitter.MessageReceived:
				if ev.Message.From.Username != self.Username {
					fmt.Printf("\n[%s]: %s\n...", ev.Message.From.Username, ev.Message.Text)
				}
			case *gitter.GitterConnectionClosed:
				fmt.Printf("!!Gitter Connection Closed!!")
				panic("!!Gitter Connection Closed!!")
			}

		case send := <-sender:
			client.SendMessage(CURRENT_ROOM_ID, send)
		}
	}
}

func get_input(receiver chan string) {
	// Get user input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("...")
		receiver <- input
	}
}

func get_env(env string) string {
	for _, element := range os.Environ() {
        key_value := strings.Split(element, "=")
		if (key_value[0] == env) {
			return key_value[1]
		}
	}
	return ""
}

func panic_err(err error) {
	if err != nil{
		panic(err)
	}
}
