// MCRCON is package for interacting with a minecraft RCON server. It provides all the necessary tools to set up a connection,
// execute commands, and recieve responses.
package mcrcon

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

/*
MCRCONPacket defines the structure of a single packet either to send to or recieve from the RCON server.

From wiki.vg/Rcon:
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

  Packet Types

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

// MCRCONClient represents a connection to a single RCON server.
// MCRCONCLient is fully synchronized and may be shared between multiple goroutines safely.
type MCRCONClient struct {
	m         sync.Mutex
	Connected bool
	conn      net.Conn
	rw        *MCRCONReaderWriter
}

// MCRCONReaderWriter is a convenience container to hold both the input and output buffers and let them be passed around easily.
type MCRCONReaderWriter struct {
	*bufio.Reader
	*bufio.Writer
}

// NewClient takes an address and a password, and attempts to set up a TCP connection to an RCON server at the given address,
// using the given password.
func NewClient(addr string, port int, passwd string) (*MCRCONClient, error) {
	nClient := MCRCONClient{}
	nClient.Connected = false

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}
	nClient.conn = conn

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	nRW := MCRCONReaderWriter{r, w}
	nClient.rw = &nRW

	// Make a pseudo-random session ID
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	openPkt := nClient.buildPacket(rnd.Int(), 3, passwd)
	nClient.writePacket(openPkt)
	authPkt, err := nClient.readPacket()

	if err != nil {
		return nil, err
	}

	// We're only connected if it returns a request type of 2.
	if authPkt.reqType == 2 {
		nClient.Connected = true
	} else {
		err := fmt.Errorf("Auth packet returned wrong type, not connected.")
		return nil, err
	}

	return &nClient, nil
}

// Close terminates RCON connection and sets the connected flag to false.
func (client *MCRCONClient) Close() {
	client.m.Lock()
	client.Connected = false
	client.conn.Close()
	client.m.Unlock()
}

// SendCommand takes a text string, executes the command on the connected client, and returns the text response
func (client *MCRCONClient) SendCommand(payload string) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	client.m.Lock()
	if client.Connected == false {
		return "", fmt.Errorf("Client not connected.")
	}
	getUserPkt := client.buildPacket(rnd.Int(), 2, payload)
	client.writePacket(getUserPkt)
	rUserPkt, rUserErr := client.readPacket()
	client.m.Unlock()

	if rUserErr != nil {
		return "", rUserErr
	} else {
		return rUserPkt.Payload, nil
	}
}

// Decode reads the RCON object's buffer and returns an MCRCONPacket representing binary data in the buffer.
// NOTE: TODO: This currently doesn't support multi-packet responses, and would behave unpredictably if one is encountered.
func (client *MCRCONClient) readPacket() (*MCRCONPacket, error) {
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

// Encode takes a MCRCONPacket and writes it into the client object's buffer in the correct binary format.
// This function assumes that the Packet is complete and correct, and just writes the results out.
func (client *MCRCONClient) writePacket(pkt *MCRCONPacket) error {
	binary.Write(client.rw, binary.LittleEndian, pkt.length)
	binary.Write(client.rw, binary.LittleEndian, pkt.reqID)
	binary.Write(client.rw, binary.LittleEndian, pkt.reqType)
	binary.Write(client.rw, binary.LittleEndian, []byte(pkt.Payload))
	binary.Write(client.rw, binary.LittleEndian, pkt.nullPad)
	return client.rw.Flush()
}

func (client *MCRCONClient) buildPacket(id int, tp int, payload string) *MCRCONPacket {
	// Build constructs an MCRCONPacket from raw information.
	newPkt := MCRCONPacket{}
	newPkt.length = int32(10 + len(payload))
	newPkt.reqID = int32(id)
	newPkt.reqType = int32(tp)
	newPkt.Payload = payload
	newPkt.nullPad[0] = 0
	newPkt.nullPad[1] = 0
	return &newPkt
}
