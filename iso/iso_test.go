package iso

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	var msg Message
	var raw = []byte("00912200302000000201000010070000000000000500001234567890211234030Hello world!, this is iso 8583")

	err := msg.Load(raw, true)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	fmt.Println(string(b))

	Equal(t, msg.ResponseCode, "1234")
	Equal(t, msg.ResponseMessage, "Hello world!, this is iso 8583")

	// py bitmap: 3020000002010000
	// go bitmap: 3020000002010000
}

func TestLongLoad(t *testing.T) {
	var msg Message
	var raw = []byte("02332200bf38404109e200080000000013000000100700000000000000000000000000000000000000000003281337210000000050573914051520190328133721032800470020040102905739140515001628112072271    005.000001     0000000006300650168888802305809876006BPJSAD")

	err := msg.Load(raw, true)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func Equal(t *testing.T, a string, b string) {
	if a != b {
		t.Errorf("%s tidak sama dengan %s", a, b)
	}
}

func TestBytes(t *testing.T) {
	var msg Message

	msg.MTI = "2200"
	msg.ProcessingCode = "100700"
	msg.Amount = "50000"
	msg.Stan = "123456789021"
	msg.ResponseCode = "1234"
	msg.ResponseMessage = "Hello world!, this is iso 8583"

	b, err := msg.Bytes(true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Raw: %v\n", string(b))
	// py: 00912200302000000201000010070000000000000500001234567890211234030Hello world!, this is iso 8583
	// go: 00912200302000000201000010070000000000000500001234567890211234030Hello world!, this is iso 8583
	// go: 0091~2200~3020000002010000~10070000000000000500001234567890211234030Hello world!, this is iso 8583

}

func TestIso(t *testing.T) {
	var msg Message

	msg.Amount = "25000"
	msg.MTI = "2200"
	msg.Stan = "123456789"

	r := reflect.TypeOf(msg)

	a, ok := r.FieldByName("Amount")
	if !ok {
		t.Error("Amount tidak ditemukan")
	}

	_, ok = a.Tag.Lookup("index")
	if !ok {
		t.Error("Nama not found")
	}

	fmt.Printf("Struct has %d member\n", r.NumField())

	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		size, _ := f.Tag.Lookup("size")
		v := reflect.ValueOf(msg).Field(i)
		// v.SetString("A")

		fmt.Printf("index %d, is %s, size: %s, value: %s\n", i, f.Name, size, v.String())
	}

}

func TestLoadEmpty(t *testing.T) {
	var data = `2200BFB8404109C30008000000001300000020070000000000000200000000000000000000000025000409162452000000000000000140498712912520200409162452031300170020040002029873837193000628112200669    352014091188404144451CA01510632548725106325487200E14C7BBDE23048858E3A9267529DA1930SYM21SB0189261130A176256B00C12CTest PLN Prepaid         R1  00000045020000180000330** INFO PEMBELIAN PLN PREPAID **
	--------------------------------
	NO METER:            51063254872
	IDPEL:              510632548720
	NAMA:  Test PLN Prepaid         
	TARIF/DAYA:             R1  /450
	--------------------------------
	****************************
	RP BAYAR : RP                  0
	--------------------------------
	006112233063012121151063254872006301212`

	var msg Message
	// var raw = []byte(data)
	raw := make([]byte, len(data))

	err := msg.Load(raw, false)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	fmt.Println(string(b))

}
