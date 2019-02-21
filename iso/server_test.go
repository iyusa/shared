package iso

import (
	"fmt"
	"testing"

	"../tool"
)

func TestServer(t *testing.T) {
	var iso Message
	iso.MTI = "2200"
	iso.ProcessingCode = PcInquiry
	iso.ResponseCode = RcFail
	iso.ResponseMessage = "This is from client"

	if err := iso.Execute("localhost", 5000); err != nil {
		t.Error(err)
	}

	Equal(t, iso.ResponseCode, RcSuccess)
	fmt.Println(tool.AsJSON(iso))
}
