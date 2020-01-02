package tool

import (
	"fmt"
	"strconv"
	"testing"
)

func TestReportString(t *testing.T) {
	rep := NewReportString(true)

	rep.AddCenter("STRUK PEMBAYARAN BPJS KESEHATAN")
	rep.AddCenter("Test Center")
	rep.Add("-")

	rep.AddKV("ID Transaksi  :", "123456")
	rep.AddKV("Tanggal Lunas :", "25-JAN-2016")
	rep.AddKV("No VA         :", "4402123")
	rep.AddKV("Nama Peserta  :", "Iyus Heraclius")
	rep.AddKV("Periode       :", "JAN04")
	rep.AddKV("Jumlah Tagihan:", "Rp 20000")
	rep.AddKV("Biaya Admin   :", "Rp 1000")
	rep.AddKV("Total         :", "Rp 21000")
	rep.Add("~")

	rep.AddCenter("BPJS KESEHATAN MENYATAKAN STRUK")
	rep.AddCenter("INI SEBAGAI BUKTI PEMBAYARAN ")
	rep.AddCenter("YANG SAH")

	fmt.Println(rep.String())
}

func TestFixed(t *testing.T) {
	raw := "040420009974   0000025000PELANGG_DENPSR_040420009974   0320170100000001675000000000-0000000020170200000001675000000000-0000000020170300000001675000000000-00000000"
	fs := NewFixedString(raw)

	fs.Add("billing_id", 15)
	fs.Add("bank_admin", 9)
	fs.Add("kode_bahasa", 1)
	fs.Add("customer_name", 30)
	fs.Add("bill_repeat_count", 2)

	totalAmount := "0"

	billCount := fs.GetInt("bill_repeat_count")
	for i := 0; i < billCount; i++ {
		k := strconv.Itoa(i + 1)
		fs.Add("bill_date_"+k, 6)
		amount, _ := fs.Add("bill_amount_"+k, 12)
		fs.Add("kubikasi_"+k, 17)
		totalAmount = SumString(totalAmount, amount)
	}

	fmt.Printf("Kubikasi 1: %s\n", fs.Get("kubikasi_1"))
	fmt.Printf("Billcount: %d\n", billCount)
	fmt.Printf("Biaya Total: %s\n", totalAmount)
	// fmt.Println(fs.Get("kode_bahasa"))
	// fmt.Println(fs.Get("customer_name"))
}

func TestMD5(t *testing.T) {
	a := "iyus"
	b := CreateMD5(a)
	fmt.Println(b)
	fmt.Println("7d0ea7482f842d31aadd256539493ab0")
}

func TestWordWraps(t *testing.T) {
	raw := "Searching for a reliable open source payment gateway solution? Continued support can put you in a reactive position, rather than place you at a proactive advantage. Does that mean you should build your own payment gateway? If you have a business based on subscriptions or a recurring revenue framework, take a look at Zuora! Our SaaS platform is purpose-built for the"
	ss := WordWraps(raw, 30)
	for _, s := range ss {
		fmt.Println(s)
	}
}
