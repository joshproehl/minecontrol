package mcrcon

import (
	"bufio"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

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
type MCRCONPacket struct {
	length, reqID, reqType int32
	Payload                string
	nullPad                [2]byte
}

type MCRCONClient struct {
	Connected bool
	conn      net.Conn
	rw        *MCRCONReaderWriter
}

type MCRCONReaderWriter struct {
	*bufio.Reader
	*bufio.Writer
}

func NewClient(addr string, passwd string) *MCRCONClient {
	nClient := MCRCONClient{}
	nClient.Connected = false

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return &nClient
	}
	nClient.conn = conn

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	nRW := MCRCONReaderWriter{r, w}
	nClient.rw = &nRW

	// Make a pseudo-random session ID
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	openPkt := nClient.Build(rnd.Int(), 3, passwd)
	nClient.Encode(openPkt)
	authPkt, err := nClient.Decode()

	if err != nil {
		return &nClient
	}

	// We're only connected if it returns a request type of 2.
	if authPkt.reqType == 2 {
		nClient.Connected = true
	}

	return &nClient
}

func (client *MCRCONClient) Close() {
	client.Connected = false
	client.conn.Close()
}

func (client *MCRCONClient) Decode() (*MCRCONPacket, error) {
	pkt := MCRCONPacket{}

	if err := binary.Read(client.rw, binary.LittleEndian, &pkt.length); err != nil {
		return &pkt, err
	}

	if err := binary.Read(client.rw, binary.LittleEndian, &pkt.reqID); err != nil {
		return &pkt, err
	}

	if err := binary.Read(client.rw, binary.LittleEndian, &pkt.reqType); err != nil {
		return &pkt, err
	}

	// Now we have the details, we'll need to load the length-10 bytes. This is because length includes the 2nd and 3rd fields (ints) and the last 2 null bytes.
	bytePayload := make([]byte, (pkt.length - 10))
	if err := binary.Read(client.rw, binary.LittleEndian, &bytePayload); err != nil {
		return &pkt, err
	}
	pkt.Payload = string(bytePayload)

	// Finally, read the last two bytes to make sure the pad is there. If these are not both NULL something went wrong.
	if err := binary.Read(client.rw, binary.LittleEndian, &pkt.nullPad); err != nil {
		return &pkt, err
	}
	// TODO: If the nullPad isn't actually NULLNULL we need to return an error.

	return &pkt, nil
}

// This function assumes that the Packet is complete and correct, and just writes the results out.
func (client *MCRCONClient) Encode(pkt *MCRCONPacket) error {
	binary.Write(client.rw, binary.LittleEndian, pkt.length)
	binary.Write(client.rw, binary.LittleEndian, pkt.reqID)
	binary.Write(client.rw, binary.LittleEndian, pkt.reqType)
	binary.Write(client.rw, binary.LittleEndian, []byte(pkt.Payload))
	binary.Write(client.rw, binary.LittleEndian, pkt.nullPad)
	return client.rw.Flush()
}

func (p *MCRCONClient) Build(id int, tp int, payload string) *MCRCONPacket {
	newPkt := MCRCONPacket{}
	newPkt.length = int32(10 + len(payload))
	newPkt.reqID = int32(id)
	newPkt.reqType = int32(tp)
	newPkt.Payload = payload
	newPkt.nullPad[0] = 0
	newPkt.nullPad[1] = 0
	return &newPkt
}

// SendCommand takes a text string, executes the command on the connected client, and returns the text response
func (client *MCRCONClient) SendCommand(payload string) (string, error) {
	if client.Connected == false {
		return "", fmt.Errorf("Client not connected.")
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	getUserPkt := client.Build(rnd.Int(), 2, payload)
	client.Encode(getUserPkt)
	rUserPkt, rUserErr := client.Decode()

	if rUserErr != nil {
		return "", rUserErr
	} else {
		return rUserPkt.Payload, nil
	}
}
