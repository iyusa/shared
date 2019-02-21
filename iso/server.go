package iso

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

// TransactionHandler interfave
type TransactionHandler interface {
	ExecuteTransaction(msg *MessageExtended) (string, error)
}

// TCPServer server handler
type TCPServer struct {
	// Config  *ini.File
	Handler TransactionHandler
}

// Serve server litener @ port
func (s *TCPServer) Serve(port int) error {
	// create tcp listener
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
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

	var ussi UssiIso
	ussi.Initialize()

	msg := UssiMessage("2210", ussi)

	if status, err := s.parseMessage(conn, msg); err != nil {
		// rc := msg.Data.(UssiIso).ResponseCode.Value
		// log.Printf("Response code %v, Error:%s \n", rc, err.Error())
		msg.writeError(conn, status, err)
		return
	}

	// Send a response back to person contacting us.
	msg.write(conn)
}

func (s *TCPServer) parseMessage(conn net.Conn, msg *MessageExtended) (string, error) {
	if s.Handler == nil {
		return RcFail, errors.New("Handler is empty")
	}

	// get first 4 bytes as length
	lenbuf := make([]byte, 4)
	reqLen, err := conn.Read(lenbuf)
	if err != nil || reqLen != 4 {
		return RcFail, err
	}

	dataLen, err := strconv.Atoi(string(lenbuf))
	if err != nil {
		return RcFail, err
	}

	// Make a buffer to hold incoming data.
	rawIso := make([]byte, dataLen)

	// Read the incoming connection into the buffer.
	reqLen, err = conn.Read(rawIso)
	if err != nil {
		return RcFail, err
	}

	// load rawIso into UssiIso
	if err := msg.Load(rawIso); err != nil {
		return RcFail, err
	}

	ussi := msg.Data.(UssiIso)
	if ussi.ProcessingCode.IsEmpty() {
		return RcFail, errors.New("Processing code is empty")
	}

	if status, err := s.Handler.ExecuteTransaction(msg); err != nil {
		return status, err
	}

	return RcSuccess, nil
}
