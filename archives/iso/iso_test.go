package iso

import (
	"fmt"
	"testing"
)

// Ussi untuk iso ussi
type Ussi struct {
	Amount *Numeric      `field:"4" length:"16"`
	Stan   *Alphanumeric `field:"11" length:"12"`
	Code   *Numeric      `field:"39" length:"4"`
	Info   *Lllvar       `field:"48" length:"999"`
	// Hp     *iso.Llvar        `field:"103" length:"28"`
}

func TestIso(t *testing.T) {
	testLoad := true

	ussi := &Ussi{
		Amount: NewNumeric("50000"),
		Stan:   NewAlphanumeric("123456789021"),
		Code:   NewNumeric("1234"),
		Info:   NewLllvar([]byte("Hello world!, this is iso 8583")),
	}

	msg := NewMessageExtended("2200", ASCII, false, true, ussi)

	if testLoad {
		raw := []byte("2200102000000201000000000000000500001234567890211234030Hello world!, this is iso 8583")
		//             2200102000000201000000000000000500001234567890211234030Hello world!, this is iso 8583
		if err := msg.Load(raw); err != nil {
			t.Error(err)
		}
	}

	b, err := msg.Bytes()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

// 2200102000000201000000000000000500001234567890211234030Hello world!, this is iso 8583	: go
// 220010200000000100000000000000050000123456789021030Hello world!, this is iso 8583		: py
