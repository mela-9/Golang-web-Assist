package reportcontroller

import (
	"Assist/config"
	"Assist/models"
	"net/http"
	"text/template"
)

// Struktur Data Laporan
type ReportData struct {
	Transactions []models.Transaction
	TotalIncome  float64
	TotalExpense float64
	Balance      float64
}

// Fungsi untuk menampilkan laporan keuangan
func ShowReport(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Ambil semua transaksi
	var transactions []models.Transaction
	err := db.Find(&transactions).Error
	if err != nil {
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}

	// Hitung total pemasukan dan pengeluaran
	var totalIncome, totalExpense float64
	for _, trx := range transactions {
		if trx.Jumlah > 0 {
			totalIncome += trx.Jumlah
		} else {
			totalExpense += trx.Jumlah
		}
	}

	// Hitung saldo
	balance := totalIncome + totalExpense

	// Kirim data ke template
	data := ReportData{
		Transactions: transactions,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      balance,
	}

	tmpl, err := template.ParseFiles("views/report/index.html")
	if err != nil {
		http.Error(w, "Gagal memuat template laporan", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
