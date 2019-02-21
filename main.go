package main

import (
	"errors"
	"fmt"

	"./iso"
	"./tool"
)

type handler struct{}

func (h *handler) ExecuteTransaction(msg *iso.Message) (string, error) {
	fmt.Printf("Receiving iso message: [%s] \n", tool.AsJSON(msg))

	if msg.ResponseCode == iso.RcFail {
		msg.ResponseCode = iso.RcSuccess
		msg.ResponseMessage = "Transakse berhasil"
		return iso.RcSuccess, nil
	}

	return iso.RcFail, errors.New("Sengaja")
}

func main() {
	h := &handler{}
	var server iso.TCPServer
	server.Handler = h
	fmt.Println("Starting server @ 5000 ...")
	server.Serve(5000)
}
