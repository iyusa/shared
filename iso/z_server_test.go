package iso

import (
	"errors"
	"fmt"
	"testing"

)

type handler struct{}

func (h *handler) ExecuteTransaction(msg *Message) (string, error) {
	// fmt.Printf("Receiving iso message: [%s] \n", tool.AsJSON(msg))

	if msg.ResponseCode == RcFail {
		msg.ResponseCode = RcSuccess
		msg.ResponseMessage = "Transaksi berhasil"
		msg.Amount = "40000"
		return RcSuccess, nil
	}

	return RcFail, errors.New("Sengaja")
}

func OffTestServer(t *testing.T) {
	var server TCPServer
	server.Handler = &handler{}
	fmt.Println("Starting server @ 5000 ...")
	server.Serve(":", 5000)
}

func TestIsoString(t *testing.T) {
	var iso Message
	iso.MTI = "2200"
	iso.ProcessingCode = PcInquiry
	iso.ResponseCode = RcFail
	iso.ResponseMessage = "This is from client"
	iso.SetAmount(50000)

	fmt.Println(iso.String())
}

// func TestClient(t *testing.T) {
// 	var iso Message
// 	iso.MTI = "2200"
// 	iso.ProcessingCode = PcInquiry
// 	iso.ResponseCode = RcFail
// 	iso.ResponseMessage = "This is from client"

// 	fmt.Println("Sending Request")
// 	fmt.Println(iso.String())

// 	if err := iso.Execute("localhost", 5000); err != nil {
// 		t.Error(err)
// 	}

// 	// Equal(t, iso.ResponseCode, RcSuccess)
// 	fmt.Println(iso.String())
// }
