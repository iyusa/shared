package tool

import (
	"fmt"
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
