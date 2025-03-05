package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/Assist"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal menghubungkan ke database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database tidak merespons:", err)
	}

	fmt.Println("Berhasil terhubung ke database")
	return db
}
