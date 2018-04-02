package lookupd

import (
	"net"
	"bufio"
	"strings"
)

type LookupProtocol struct {
	ctx *Context
}

func (p *LookupProtocol) IOLoop(conn net.Conn) error {
	var err error
	var line string
	client := NewClient(conn)
	reader := bufio.NewReader(client)
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		// TODO what about these
		line = strings.TrimSpace(line)
		params := strings.Split(line, " ")

		var response []byte
		response, err = p.Exec(client, reader, params)
	}
	return  nil
}
