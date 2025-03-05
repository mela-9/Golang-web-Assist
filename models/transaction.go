package models

import (
	"Golang-web-Assist/entities"
	"database/sql"
)

func GetTransactions(db *sql.DB) ([]entities.Transaction, error) {
	rows, err := db.Query("SELECT id, amount, category, date, description FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Category, &t.Date, &t.Description); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
func InsertTransaction(db *sql.DB, transaction entities.Transaction) error {
	query := "INSERT INTO transactions (amount, category, date, description) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, transaction.Amount, transaction.Category, transaction.Date, transaction.Description)
	return err
}
