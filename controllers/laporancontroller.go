package controllers

import (
	"Golang-web-Assist/config"
	"Golang-web-Assist/entities"
	"Golang-web-Assist/models"
	"html/template"
	"log"
	"net/http"
)

// Fungsi untuk menampilkan laporan keuangan
func LaporanHandler(w http.ResponseWriter, r *http.Request) {
	// Menghubungkan ke database
	db := config.ConnectDB()
	defer db.Close()

	// Ambil data transaksi dari model dengan mengirimkan objek db
	transactions, err := models.GetTransactions(db)
	if err != nil {
		log.Println("Error retrieving transactions:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parsing template laporan
	t, err := template.ParseFiles("views/report/report.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Kirim data transaksi ke template
	data := struct {
		Title        string
		Transactions []entities.Transaction
	}{
		Title:        "Laporan Keuangan",
		Transactions: transactions,
	}

	// Eksekusi template dengan data yang sudah disiapkan
	err = t.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
