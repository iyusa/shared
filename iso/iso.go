package iso

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/iyusa/shared/tool"
)

const (
	tagIndex  string = "index"  // BIT
	tagLength string = "length" // 1, 99
	tagType   string = "type"   // A, N, AN
	tagSize   string = "size"   // FIXED, LLVAR, LLLVAR
)

const (
	sizeFixed  = "FIXED"
	sizeLLVAR  = "LLVAR"
	sizeLLLVAR = "LLLVAR"
)

const (
	typeNumeric      = "N"
	typeAlphaNumeric = "AN"
)

// Message main struct
type Message struct {
	MTI string

	ProcessingCode  string `index:"3" length:"6" type:"N" size:"FIXED"`
	Amount          string `index:"4" length:"16" type:"N" size:"FIXED"`
	Stan            string `index:"11" length:"12" type:"N" size:"FIXED"`
	TransactionTime string `index:"12" length:"14" type:"N" size:"FIXED"`
	Period          string `index:"40" length:"3" type:"N" size:"FIXED"`
	ResponseCode    string `index:"39" length:"4" type:"N" size:"FIXED"`
	Buffer          string `index:"47" length:"999" type:"AN" size:"LLLVAR"`
	ResponseMessage string `index:"48" length:"999" type:"AN" size:"LLLVAR"`
	Extra1          string `index:"61" length:"999" type:"AN" size:"LLLVAR"`
	Extra2          string `index:"62" length:"999" type:"AN" size:"LLLVAR"`
	BillerCode      string `index:"100" length:"99" type:"AN" size:"LLVAR"`
	SubscriberID    string `index:"103" length:"99" type:"AN" size:"LLVAR"`
	ProductCode     string `index:"104" length:"999" type:"AN" size:"LLLVAR"`
}

// Bytes create []byte representation
func (m *Message) Bytes(withLength bool) ([]byte, error) {
	data := ""
	bits := []int{0, 1, 128, 64, 32, 16, 8, 4, 2}
	bitmap := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	r := reflect.TypeOf(*m)
	v := reflect.ValueOf(m).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)

		idx, ok := f.Tag.Lookup(tagIndex)
		if ok {
			index, _ := strconv.Atoi(idx)
			val := v.Field(i).String()

			value := getFieldValue(f, val) // masih mengandung len untuk LLVAR dan LLLVAR
			if value != "" {
				data += value

				if index > 64 {
					bitmap[0] = bitmap[0] | bits[2]
				}

				pos := index / 8
				if index%8 == 0 {
					pos = (index / 8) - 1
				}

				// pos := (index % 8) == 0 ? (index / 8) - 1 : index / 8
				bitmap[pos] = bitmap[pos] | bits[(index%8)+1]
			}
		}
	}

	// build bitmap hex
	bitmapHex := ""
	for c := 0; c < 16; c++ {
		tm := fmt.Sprintf("%X", bitmap[c])
		if len(tm) < 2 {
			tm = "0" + tm
		}
		bitmapHex += tm

		if (bitmap[0]&128) != 128 && (c == 7) {
			break
		}
	}

	s := m.MTI + bitmapHex + data
	if withLength {
		n := len(s)
		ls := fmt.Sprintf("%04d", n)
		s = ls + s
	}
	ret := []byte(s)
	return ret, nil
}

// value sudah disesuaikan dengan type, size
func getFieldValue(f reflect.StructField, val string) string {
	flength, _ := f.Tag.Lookup(tagLength)
	typ, _ := f.Tag.Lookup(tagType)
	siz, _ := f.Tag.Lookup(tagSize)

	length, _ := strconv.Atoi(flength)

	if val == "" {
		return val
	}

	if len(val) > length {
		val = val[:length]
	}

	if siz == sizeFixed {
		if typ == typeNumeric {
			return Left(val, length, "0")
		}
		return Right(val, length, " ")
	}

	if siz == sizeLLVAR {
		ln := len(val)
		s := fmt.Sprintf("%02d", ln)
		return s + val
	}

	if siz == sizeLLLVAR {
		ln := len(val)
		s := fmt.Sprintf("%03d", ln)
		return s + val
	}

	return ""
}

// Load from stream
func (m *Message) Load(source []byte, hasLength bool) error {
	if len(source) < 24 {
		return errors.New("Invalid raw message (too small)")
	}

	var offset int
	if hasLength {
		offset = 4
	}

	// read mti
	b := source[offset : offset+4]
	offset += 4
	m.MTI = string(b)

	// read bitmap
	bitmapHex, bitmaps := buildBitmap(source, offset)
	offset += len(bitmapHex)

	return m.buildValues(source, bitmapHex, bitmaps, offset)
}

func buildBitmap(source []byte, offset int) (bitmapHex string, bitmaps []int) {
	//bitmaps = make([]int, 16)
	bitmaps = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bitmapHex = ""

	top := offset + 32
	mid := offset + 14
	k := 0

	s := source[offset : offset+2]
	first, _ := strconv.ParseInt(string(s), 16, 0)
	noext := (first & 128) != 128

	for offset < top { // nanti check, karena asalnya while
		tmp := source[offset : offset+2]
		hexStr := string(tmp)
		bitmapHex += hexStr
		x, _ := strconv.ParseInt(hexStr, 16, 0)
		bitmaps[k] = int(x)

		if noext && (offset == mid) {
			break
		}

		k++
		offset += 2
	}
	bitmapHex = strings.ToUpper(bitmapHex)
	return
}

func (m *Message) buildValues(source []byte, bitmapHex string, bitmaps []int, offset int) error {
	bitmapValues := make([]bool, 128)
	bits := []int{0, 1, 128, 64, 32, 16, 8, 4, 2}

	// flag yang akan diisi
	for i := 0; i < 16; i++ {
		for j := 1; j < 9; j++ {
			if (bitmaps[i] & bits[j]) == bits[j] {
				if j == 1 {
					bitmapValues[(i+1)*8] = true
				} else if j != 2 || i != 0 {
					bitmapValues[i*8+j-1] = true
				}
			}
		}
	}

	// offset += 2

	r := reflect.TypeOf(*m)
	v := reflect.ValueOf(m).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)

		idx, ok := f.Tag.Lookup(tagIndex)
		if ok {
			fval := v.Field(i)
			if !fval.CanSet() {
				return errors.New("Field tidak bisa di set")
			}

			index, _ := strconv.Atoi(idx)
			if bitmapValues[index] {
				flength, _ := f.Tag.Lookup(tagLength)
				siz, _ := f.Tag.Lookup(tagSize)
				length, _ := strconv.Atoi(flength)

				if siz == sizeLLVAR {
					valsize := btoi(source[offset : offset+2])
					if valsize > length || valsize == 0 {
						return errors.New("Ukuran LLVar terlalu panjang atau tidak diisi")
					}

					value := string(source[offset+2 : offset+2+valsize])
					fval.SetString(value)
					offset += valsize + 2
				} else if siz == sizeLLLVAR {
					valsize := btoi(source[offset : offset+3])
					if valsize > length || valsize == 0 {
						return errors.New("Ukuran LLLVar terlalu panjang atao kosong")
					}

					value := string(source[offset+3 : offset+3+valsize])
					fval.SetString(value)
					offset += valsize + 3
				} else { // berarti fixed
					value := string(source[offset : offset+length])
					fval.SetString(value)
					offset += length
				}

			}
		}
	}

	return nil
}

// Write to connection()
func (m *Message) Write(conn net.Conn) error {
	buf, err := m.Bytes(true)
	if err != nil {
		return err
	}
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

// WriteError default error to connection
func (m *Message) WriteError(conn net.Conn, status string, err error) error {
	m.ResponseCode = status
	m.ResponseMessage = err.Error()

	return m.Write(conn)
}

// Execute send iso to host
func (m *Message) Execute(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// send data
	if err := m.Write(conn); err != nil {
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

	fmt.Printf("Receiving raw iso %s\n", string(rawIso))

	// load rawIso into msg.Data / UssiIso
	if err := m.Load(rawIso, false); err != nil {
		return err
	}

	return nil
}

// SetAmount as integer
func (m *Message) SetAmount(amount int) {
	m.Amount = string(amount)
}

// GetAmount as integer
func (m *Message) GetAmount() int {
	return tool.StrToInt(m.Amount, 0)
}

// https://github.com/willf/pad/blob/master/pad.go
func times(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

// Left left-pads the string with pad up to len runes
// len may be exceeded if
func Left(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

// Right right-pads the string with pad up to len runes
func Right(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

func btoi(b []byte) int {
	s := string(b)
	a, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return a
}
