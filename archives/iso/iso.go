package iso

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// TransactionHandler interfave
type TransactionHandler interface {
	ExecuteTransaction(msg *MessageExtended) (string, error)
}

const (
	PcInquiry  = "100700"
	PcPayment  = "200700"
	PcPurchase = "300700"
	PcAdvice   = "301700"
)

const (
	RcSuccess = "0000"
	RcPending = "0068"
	RcFail    = "1000"
)

func byt(source string) []byte {
	return []byte(source)
}

// UssiIso untuk iso ussi
type UssiIso struct {
	ProcessingCode  *Numeric      `field:"3" length:"6"`
	Amount          *Numeric      `field:"4" length:"16"`
	Stan            *Alphanumeric `field:"11" length:"12"`
	TransactionTime *Numeric      `field:"12" length:"12"`
	Period          *Numeric      `field:"40" length:"3"`
	ResponseCode    *Numeric      `field:"39" length:"4"`
	Buffer          *Lllvar       `field:"47" length:"999"`
	ResponseMessage *Lllvar       `field:"48" length:"999"`
	Extra1          *Lllvar       `field:"61" length:"999"`
	Extra2          *Lllvar       `field:"62" length:"999"`
	BillerCode      *Llvar        `field:"100" length:"99"`
	SubscriberID    *Llvar        `field:"103" length:"99"`
	ProductCode     *Lllvar       `field:"104" length:"999"`
}

// Initialize internal fields
func (u *UssiIso) Initialize() {
	u.ProcessingCode = NewNumeric("")
	u.Amount = NewNumeric("")
	u.Stan = NewAlphanumeric("")
	u.TransactionTime = NewNumeric("")
	u.Period = NewNumeric("")
	u.ResponseCode = NewNumeric("")
	u.Buffer = NewLllvar([]byte(""))
	u.ResponseMessage = NewLllvar(byt(""))
	u.Extra1 = NewLllvar([]byte(""))
	u.Extra2 = NewLllvar([]byte(""))
	u.BillerCode = NewLlvar([]byte(""))
	u.SubscriberID = NewLlvar([]byte(""))
	u.ProductCode = NewLllvar([]byte(""))
}

func (u *UssiIso) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		return "~"
	}
	return string(b)
}

// UssiMessage shortcut to NewMessageExtended
func UssiMessage(mti string, data interface{}) *MessageExtended {
	return NewMessageExtended(mti, ASCII, true, true, data)
}

// func (m *MessageExtended) parse(source string) error {
// 	raw := source[4:]
// 	err := m.Load([]byte(raw))
// 	return err
// }

func (m *MessageExtended) toByte(withLength bool) ([]byte, error) {
	buf, err := m.Bytes()
	if err != nil {
		return nil, err
	}

	if withLength {
		s := fmt.Sprintf("%04d", len(buf))
		result := append([]byte(s), buf...)
		return result, nil
	}
	return buf, nil
}

// Execute send iso to host
func (m *MessageExtended) Execute(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// send data
	if err := m.write(conn); err != nil {
		return err
	}

	// read response
	// get first 4 bytes as length
	lenbuf := make([]byte, 4)
	reqLen, err := conn.Read(lenbuf)
	if err != nil || reqLen != 4 {
		return err
	}

	dataLen, err := strconv.Atoi(string(lenbuf))
	if err != nil {
		return err
	}

	// Make a buffer to hold incoming data.
	rawIso := make([]byte, dataLen)

	// Read the incoming connection into the buffer.
	reqLen, err = conn.Read(rawIso)
	if err != nil {
		return err
	}

	// load rawIso into msg.Data / UssiIso
	if err := m.Load(rawIso); err != nil {
		return err
	}

	return nil
}

func (m *MessageExtended) write(conn net.Conn) error {
	buf, err := m.toByte(true)
	if err != nil {
		return err
	}
	conn.Write(buf)
	// log.Printf("Buffer:[%s]\n", string(buf))
	return nil
}

func (m *MessageExtended) writeError(conn net.Conn, status string, err error) {
	iso := m.Data.(UssiIso)

	iso.ResponseCode.Value = status
	iso.ResponseMessage.Value = byt(err.Error())

	m.write(conn)
}
