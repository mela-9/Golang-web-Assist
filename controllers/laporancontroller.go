package controllers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"Golang-web-Assist/config" // Pastikan import config database
)

// Struct untuk transaksi
type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Date        string
	Description string
}

// Struct untuk menyimpan data laporan keuangan
type ReportData struct {
	Transactions  []Transaction
	TotalIncome   float64
	TotalExpense  float64
	Balance       float64
	TargetSavings float64
}

// Handler untuk halaman laporan
func ShowReport(w http.ResponseWriter, r *http.Request) {
	// Koneksi ke database
	db := config.ConnectDB()
	defer db.Close()

	// Ambil saldo awal dari tabel initial_balance
	var initialBalance float64
	err := db.QueryRow("SELECT amount FROM initial_balance ORDER BY id DESC LIMIT 1").Scan(&initialBalance)
	if err != nil {
		log.Println("⚠️ Tidak ada saldo awal yang tersimpan, default ke 0")
		initialBalance = 0 // Default jika belum ada saldo awal
	}

	// Query ambil semua transaksi
	rows, err := db.Query("SELECT id, amount, category, date, description FROM transactions ORDER BY date DESC")
	if err != nil {
		log.Println("❌ Gagal mengambil data transaksi:", err)
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	var totalIncome, totalExpense float64

	// Looping ambil data
	for rows.Next() {
		var trx Transaction
		var date string

		if err := rows.Scan(&trx.ID, &trx.Amount, &trx.Category, &date, &trx.Description); err != nil {
			log.Println("❌ Gagal membaca data transaksi:", err)
			continue
		}

		trx.Date = date // Simpan langsung sebagai string
		// Hitung total pemasukan dan pengeluaran berdasarkan kategori
		if trx.Amount > 0 {
			totalIncome += trx.Amount
		} else {
			totalExpense += -trx.Amount
		}

		transactions = append(transactions, trx)
	}

	for i := range transactions {
		transactions[i].ID = i + 1 // Nomor ulang ID berdasarkan indeks array
	}

	// Hitung saldo akhir
	balance := initialBalance + totalIncome - totalExpense

	// Query untuk mengambil target tabungan terbaru
	var targetSavings float64
	err = db.QueryRow("SELECT target_amount FROM savings ORDER BY id DESC LIMIT 1").Scan(&targetSavings)
	if err != nil {
		log.Println("⚠️ Tidak ada target tabungan yang tersimpan, default 0")
		targetSavings = 0 // Set default jika tidak ada data
	}

	// Buat struct yang menggabungkan transaksi dan target tabungan
	reportData := ReportData{
		Transactions:  transactions,
		TotalIncome:   totalIncome,
		TotalExpense:  totalExpense,
		Balance:       balance,
		TargetSavings: targetSavings,
	}

	// Debugging output
	fmt.Println("✅ Saldo Awal:", initialBalance)
	fmt.Println("✅ Total Pemasukan:", totalIncome)
	fmt.Println("✅ Total Pengeluaran:", totalExpense)
	fmt.Println("✅ Saldo Akhir:", balance)
	fmt.Println("✅ Data transaksi:", transactions)

	// Load template
	tmpl, err := template.ParseFiles("views/report/report.html")
	if err != nil {
		log.Println("❌ Gagal memuat template:", err)
		http.Error(w, "Gagal memuat halaman", http.StatusInternalServerError)
		return
	}
	fmt.Println("✅ Data untuk laporan:", reportData)
	tmpl.Execute(w, reportData)
}
