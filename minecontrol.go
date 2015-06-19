package main

import (
  "fmt"
  "bufio"
  "io"
  "encoding/binary"
)


// A passthrough type allowing us a customer "parser class"
type RCONParser struct {
  *bufio.Reader
}

// Get an instance of our "parser class"
func NewRCONParser(r io.Reader) *RCONParser {
  return &RCONParser{
    Reader: bufio.NewReader(r),
  }
}

/*
From wiki.vg/Rcon

Packet Format

Integers are little-endian, in contrast with the Beta protocol.
Responses are sent back with the same Request ID that you send.
In the event of an auth failure (i.e. your login is incorrect, or you're trying to send commands without first logging in), request ID will be set to -1

Field name  Field Type  Notes
Length      int         Length of remainder of packet
Request ID  int         Client-generated ID
Type        int         3 for login, 2 to run a command, 0 for a multi-packet response
Payload     byte[]      ASCII text
2-byte pad  byte, byte  Two null bytes

Packets

3: Login
  Outgoing payload: password.
  If the server returns a packet with the same request ID, auth was successful (note: packet type is 2, not 3). If you get an request ID of -1, auth failed (wrong password).
2: Command
  Outgoing payload should be the command to run, e.g. time set 0
0: Command response
  Incoming payload is the output of the command, though many commands return nothing, and there's no way of detecting unknown commands.
  The output of the command may be split over multiple packets, each containing 4096 bytes (less for the last packet). Each packet contains part of the payload (and the two-byte padding). The last packet sent is the end of the output.

Maximum request length: 1460 (giving a max payload length of 1446)
Code exists in the notchian server to split large responses (>4096 bytes) into multiple smaller packets. However, the code that actually encodes each packet expects a max length of 1248, giving a max response payload length of 1234 bytes.
*/
type RCONPacket struct {
  length, reqID, reqType int
  payload []byte
  nullPad [2]byte
}

func (reader *RCONParser) Decode() (*RCONPacket, error) {
  pkt := RCONPacket{}

  if err := binary.Read(reader, binary.LittleEndian, &pkt.length); err != nil {
    return &pkt, err
  }

  if err := binary.Read(reader, binary.LittleEndian, &pkt.reqID); err != nil {
    return &pkt, err
  }

  if err := binary.Read(reader, binary.LittleEndian, &pkt.reqType); err != nil {
    return &pkt, err
  }

  // Now we have the details, we'll need to load the length-10 bytes. This is because length includes the 2nd and 3rd fields (ints) and the last 2 null bytes.
  pkt.payload =  make([]byte, pkt.length)
  if err := binary.Read(reader, binary.LittleEndian, &pkt.payload); err != nil {
    return &pkt, err
  }


  // Finally, read the last two bytes to make sure the pad is there. If these are not both NULL something went wrong.
  if err := binary.Read(reader, binary.LittleEndian, &pkt.nullPad); err != nil {
    return &pkt, err
  }
  // TODO: If the nullPad isn't actually NULLNULL we need to return an error.

  return &pkt, nil
}


func main() {
  fmt.Println("Minecontrol running...")
}


