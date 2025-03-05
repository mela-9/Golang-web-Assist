package controllers

import (
	"database/sql"
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

// Handler untuk halaman laporan
func ShowReport(w http.ResponseWriter, r *http.Request) {
	// Koneksi ke database
	db := config.ConnectDB()
	defer db.Close()

	// Query ambil semua transaksi
	rows, err := db.Query("SELECT id, amount, category, date, description FROM transactions ORDER BY date DESC")
	if err != nil {
		log.Println("❌ Gagal mengambil data transaksi:", err)
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction

	// Looping ambil data
	for rows.Next() {
		var trx Transaction
		var date sql.NullTime

		if err := rows.Scan(&trx.ID, &trx.Amount, &trx.Category, &date, &trx.Description); err != nil {
			log.Println("❌ Gagal membaca data transaksi:", err)
			continue
		}

		// Konversi tanggal dari `time.Time` ke `string`
		if date.Valid {
			trx.Date = date.Time.Format("02-01-2006")
		} else {
			trx.Date = "Tidak diketahui"
		}

		transactions = append(transactions, trx)
	}

	// Cek apakah transaksi berhasil diambil
	fmt.Println("✅ Data transaksi:", transactions)

	// Load template
	tmpl, err := template.ParseFiles("views/report/report.html")
	if err != nil {
		log.Println("❌ Gagal memuat template:", err)
		http.Error(w, "Gagal memuat halaman", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, transactions)
}
