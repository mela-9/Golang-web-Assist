package reportcontroller

import (
	"Golang-web-Assist/config"
	"Golang-web-Assist/entities"
	"log"
	"net/http"
	"text/template"
)

// Struktur Data Laporan
type ReportData struct {
	Transactions []entities.Transaction
	TotalIncome  float64
	TotalExpense float64
	Balance      float64
}

// Fungsi untuk menampilkan laporan keuangan
func ShowReport(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Ambil semua transaksi dari database menggunakan Query
	rows, err := db.Query("SELECT id, amount, category, date, description FROM transactions")
	if err != nil {
		log.Println("Error saat mengambil data transaksi:", err)
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []entities.Transaction
	var totalIncome, totalExpense float64

	// Loop melalui hasil query
	for rows.Next() {
		var trx entities.Transaction
		err := rows.Scan(&trx.ID, &trx.Amount, &trx.Category, &trx.Date, &trx.Description)
		if err != nil {
			log.Println("Error saat membaca data transaksi:", err)
			http.Error(w, "Gagal membaca data transaksi", http.StatusInternalServerError)
			return
		}

		// Hitung total pemasukan & pengeluaran
		if trx.Amount > 0 {
			totalIncome += trx.Amount
		} else {
			totalExpense += trx.Amount
		}

		transactions = append(transactions, trx)
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

	// Parsing template laporan
	tmpl, err := template.ParseFiles("views/report/index.html")
	if err != nil {
		log.Println("Error saat memuat template:", err)
		http.Error(w, "Gagal memuat template laporan", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
