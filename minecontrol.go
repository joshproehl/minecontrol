package main

import (
  "fmt"
  "github.com/joshproehl/minecontrol-go/mcrcon"
)

func main() {
  fmt.Println("Minecontrol running...")


  client := mcrcon.NewClient("127.0.0.1:25575", "password")
  defer client.Close()


  openPkt := client.Build(666, 3, "password")
  fmt.Println("Sending packet: ", openPkt)

  client.Encode(openPkt)
  packet, err := client.Decode()

  if(err != nil) {
    fmt.Println("FATAL: ", err)
  }

  fmt.Println("Auth Result packet was", packet)

  getUserPkt := client.Build(666, 2, "/list")
  client.Encode(getUserPkt)
  rUserPkt, rUserErr := client.Decode()

  if(rUserErr != nil) {
    fmt.Println("FATAL: ", rUserErr)
  }

  fmt.Println("User result was: ", rUserPkt)
  fmt.Println(rUserPkt.Payload)

}


