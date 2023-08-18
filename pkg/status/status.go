package status

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ServerStatus struct {
	version           string
	users_online      uint32
	users_max         uint32
	allowed_bandwidth uint32 // in bits/s
}

func GetServerStatus(address string) (ServerStatus, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return ServerStatus{}, err
	}

	//                     |    PING COMMAND    |  |            MESSAGE IDENTIFIER              |
	ping_message := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	binary.BigEndian.PutUint64(ping_message[4:12], uint64(time.Now().UnixMicro()))
	conn.Write(ping_message)

	response := make([]byte, 24)
	_, err = bufio.NewReader(conn).Read(response)
	if err != nil {
		return ServerStatus{}, err
	}

	version := fmt.Sprintf("%d.%d.%d", response[1], response[2], response[3])
	// this will be correct, no need to verify :)
	// ident := binary.BigEndian.Uint64(response[4:12])
	users_online := binary.BigEndian.Uint32(response[12:16])
	users_max := binary.BigEndian.Uint32(response[16:20])
	allowed_bandwidth := binary.BigEndian.Uint32(response[20:24])

	return ServerStatus{
		version:           version,
		users_online:      users_online,
		users_max:         users_max,
		allowed_bandwidth: allowed_bandwidth,
	}, nil
}
