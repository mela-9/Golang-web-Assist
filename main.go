package main

import (
	"fmt"
	"log"
	"net/http"

	"Golang-web-Assist/config"
	"Golang-web-Assist/controllers/homecontroller"
	"Golang-web-Assist/controllers/reportcontroller"
	"Golang-web-Assist/controllers/transactioncontroller"
)

func main() {
	// Koneksi ke database
	db := config.ConnectDB()
	defer db.Close()

	// **Deklarasi router (baru)**
	mux := http.NewServeMux() // ✅ Buat router sebelum dipakai

	// Routing halaman utama (Dashboard)
	mux.HandleFunc("/", homecontroller.Dashboard)

	// Routing halaman tambah transaksi
	mux.HandleFunc("/add-transaction", transactioncontroller.TambahTransaksi)

	// Routing halaman laporan keuangan & tabungan
	mux.HandleFunc("/report", reportcontroller.ShowReport)

	// Jalankan server
	fmt.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
