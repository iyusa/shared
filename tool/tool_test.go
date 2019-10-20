package tool

import (
	"fmt"
	"strconv"
	"testing"
)

func TestReportString(t *testing.T) {
	rep := NewReportString()

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
