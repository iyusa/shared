package server

import (
	"errors"
	"fmt"
	"testing"

	"../iso"
	"../tool"
)

type handler struct{}

func (h *handler) Execute(msg *iso.Message) error {
	fmt.Printf("Receiving iso message: [%s] \n", tool.AsJSON(msg))

	if msg.ResponseCode == iso.RcFail {
		msg.ResponseCode = iso.RcSuccess
		msg.ResponseMessage = "Transaksi berhasil"
		msg.Amount = "40000"
		return nil
	}

	return errors.New("Sengaja")
}

func OffTestServer(t *testing.T) {
	var server IsoServer
	server.Handler = &handler{}
	fmt.Println("Starting server @ 5000 ...")
	server.Serve(":", 5000)
}

func TestIsoString(t *testing.T) {
	var msg iso.Message
	msg.MTI = "2200"
	msg.ProcessingCode = iso.PcInquiry
	msg.ResponseCode = iso.RcFail
	msg.ResponseMessage = "This is from client"
	msg.SetAmount(50000)

	fmt.Println(msg.String())
}

// func TestClient(t *testing.T) {
// 	var msg iso.Message
// 	msg.MTI = "2200"
// 	msg.ProcessingCode = iso.PcInquiry
// 	msg.ResponseCode = iso.RcFail
// 	msg.ResponseMessage = "This is from client"

// 	fmt.Println("Sending Request")
// 	fmt.Println(msg.String())

// 	if err := msg.Execute("localhost", 5000); err != nil {
// 		t.Error(err)
// 	}

// 	// Equal(t, iso.ResponseCode, RcSuccess)
// 	fmt.Println(msg.String())
// }
