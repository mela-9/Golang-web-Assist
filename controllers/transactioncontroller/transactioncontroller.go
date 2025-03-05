package transactioncontroller

import (
	"html/template"
	"net/http"
	"strconv"

	"Golang-web-Assist/config"
	"Golang-web-Assist/entities"
	"Golang-web-Assist/models"
)

func TambahTransaksi(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db := config.ConnectDB()
		defer db.Close()

		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		transaction := entities.Transaction{
			Amount:      amount,
			Category:    r.FormValue("category"),
			Date:        r.FormValue("date"),
			Description: r.FormValue("description"),
		}

		err := models.InsertTransaction(db, transaction)
		if err != nil {
			http.Error(w, "Gagal menyimpan transaksi", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("views/home/add_transaction.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman tambah transaksi", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Edit Transaksi
func EditTransaksi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil data dari form
	id := r.FormValue("id")
	deskripsi := r.FormValue("deskripsi")
	jumlah := r.FormValue("jumlah")
	kategori := r.FormValue("kategori")

	// Konversi jumlah ke integer
	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil {
		http.Error(w, "Jumlah harus berupa angka", http.StatusBadRequest)
		return
	}

	// Koneksi database
	db := config.ConnectDB()
	defer db.Close()

	// Update data transaksi
	_, err = db.Exec("UPDATE transactions SET deskripsi=?, jumlah=?, kategori=? WHERE id=?", deskripsi, jumlahInt, kategori, id)
	if err != nil {
		http.Error(w, "Gagal mengupdate transaksi", http.StatusInternalServerError)
		return
	}

	// Redirect ke halaman utama setelah edit
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Hapus Transaksi
func HapusTransaksi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil ID transaksi yang akan dihapus
	id := r.FormValue("id")

	// Koneksi database
	db := config.ConnectDB()
	defer db.Close()

	// Hapus transaksi berdasarkan ID
	_, err := db.Exec("DELETE FROM transactions WHERE id=?", id)
	if err != nil {
		http.Error(w, "Gagal menghapus transaksi", http.StatusInternalServerError)
		return
	}

	// Redirect ke halaman utama setelah hapus
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
