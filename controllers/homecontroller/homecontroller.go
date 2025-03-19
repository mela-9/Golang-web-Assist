package homecontroller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"Golang-web-Assist/config"
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

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Date        string
	Description string
}

func ShowReport(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Ambil data transaksi
	rows, err := db.Query("SELECT id, amount, category, date, description FROM transactions ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Gagal mengambil data transaksi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var trx Transaction
		if err := rows.Scan(&trx.ID, &trx.Amount, &trx.Category, &trx.Date, &trx.Description); err == nil {
			transactions = append(transactions, trx)
		}
	}

	// Ambil target tabungan terbaru
	var targetTabungan float64
	err = db.QueryRow("SELECT target_amount FROM savings ORDER BY id DESC LIMIT 1").Scan(&targetTabungan)
	if err != nil {
		log.Println("⚠️ Tidak ada target tabungan yang tersimpan")
		targetTabungan = 0
	}

	// Gabungkan data transaksi & target tabungan ke struct
	data := struct {
		Transactions  []Transaction
		TargetSavings float64
	}{
		Transactions:  transactions,
		TargetSavings: targetTabungan,
	}

	// Render ke template laporan.html
	tmpl, _ := template.ParseFiles("views/report/report.html")
	tmpl.Execute(w, data)
}
