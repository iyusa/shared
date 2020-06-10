package iso

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// TransactionHandler interfave. Deprecated, user server.Server instead
type TransactionHandler interface {
	ExecuteTransaction(msg *Message) (string, error)
}

// TCPServer server handler
type TCPServer struct {
	Handler TransactionHandler
}

// Serve server litener @ localhost:port
func (s *TCPServer) Serve(host string, port int) error {
	// create tcp listener
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	// l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer l.Close()

	// Listen for an incoming connection.
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}

		// Handle connections in a new goroutine.
		go s.handleRequest(conn)
	}
}

func (s *TCPServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	msg := &Message{}
	msg.MTI = "2100"

	// 1. parse iso message from connection
	if status, err := s.parseMessage(conn, msg); err != nil {
		msg.WriteError(conn, status, err)
		return
	}

	// 2. execute transaction
	if status, err := s.Handler.ExecuteTransaction(msg); err != nil {
		msg.WriteError(conn, status, err)
		return
	}

	// 3. send back iso to caller
	msg.Write(conn)

	conn.Close()
}

// parse message from connection into msg (msg already created)
func (s *TCPServer) parseMessage(conn net.Conn, msg *Message) (string, error) {
	if s.Handler == nil {
		return RcFail, errors.New("Handler is empty")
	}

	// get first 4 bytes as length
	lenbuf := make([]byte, 4)
	reqLen, err := conn.Read(lenbuf)
	if err != nil {
		return RcFail, err
	}

	if reqLen != 4 {
		return RcFail, fmt.Errorf("Invalid buffer length")
	}

	dataLen, err := strconv.Atoi(string(lenbuf))
	if err != nil {
		return RcFail, err
	}

	// Make a buffer to hold incoming data.
	rawIso := make([]byte, dataLen)

	// Read the incoming connection into the buffer.
	reader := bufio.NewReader(conn)
	reqLen, err = io.ReadFull(reader, rawIso)
	// reqLen, err = conn.Read(rawIso)
	if err != nil {
		return RcFail, err
	}

	// load rawIso into UssiIso
	if err := msg.Load(rawIso, false); err != nil {
		return RcFail, err
	}

	return RcSuccess, nil
}
