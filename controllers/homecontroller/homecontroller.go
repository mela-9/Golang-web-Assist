package homecontroller

import (
	"html/template"
	"net/http"
	"strconv"

	"Golang-web-Assist/config" // Pastikan sesuai dengan module di go.mod
	"Golang-web-Assist/entities"
	"Golang-web-Assist/models"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	transactions, err := models.GetTransactions(db)
	if err != nil {
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("views/home/dashboard.html")
	if err != nil {
		http.Error(w, "Gagal memuat dashboard", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, transactions)
}

func LaporanKeuangan(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home/report.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman laporan keuangan", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home/index.html")
	if err != nil {
		http.Error(w, "Gagal memuat template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	db := config.ConnectDB()
	defer db.Close()

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Jumlah tidak valid", http.StatusBadRequest)
		return
	}

	transaction := entities.Transaction{
		Amount:      amount,
		Category:    r.FormValue("category"),
		Date:        r.FormValue("date"),
		Description: r.FormValue("description"),
	}
	err = models.InsertTransaction(db, transaction)
	if err != nil {
		http.Error(w, "Gagal menyimpan transaksi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
