package main

import (
	"fmt"
	"github.com/joshproehl/minecontrol-go/mcrcon"
)

func main() {
	fmt.Println("Minecontrol running...")

	client := mcrcon.NewClient("127.0.0.1:25575", "password")

	if client.Connected != true {
		fmt.Println("FATAL: Client could not connect")
		return
	}

	fmt.Println("Getting player list...")

	listResponse, rUserErr := client.SendCommand("/list")

	if rUserErr != nil {
		fmt.Println("FATAL: ", rUserErr)
	}

	fmt.Println(listResponse)

	client.Close()
}
