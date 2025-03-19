package controllers

import (
	"Golang-web-Assist/config"
	"log"
	"net/http"
	"strconv"
)

// SaveSavings menyimpan target tabungan ke database
func SaveSavings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	targetStr := r.FormValue("target")
	if targetStr == "" {
		http.Error(w, "Target tidak boleh kosong", http.StatusBadRequest)
		return
	}

	target, err := strconv.ParseFloat(targetStr, 64)
	if err != nil {
		log.Println("❌ Gagal mengonversi target:", err)
		http.Error(w, "Format angka tidak valid", http.StatusBadRequest)
		return
	}

	db := config.ConnectDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO savings (target_amount) VALUES (?)")
	if err != nil {
		log.Println("❌ Gagal mempersiapkan query:", err)
		http.Error(w, "Gagal menyimpan target", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(target)
	if err != nil {
		log.Println("❌ Gagal menyimpan target tabungan:", err)
		http.Error(w, "Gagal menyimpan target", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Target tabungan berhasil disimpan:", target)
	http.Redirect(w, r, "/laporan", http.StatusSeeOther)
}
