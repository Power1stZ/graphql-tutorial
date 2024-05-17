package deposit

import (
	"fmt"
	"graphql/tutorial/pkg/db"
	"log/slog"
	"time"
)

type Transaction struct {
	TransactionID int       `json:"trans_id"`
	Amount        int       `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	FinishedAt    time.Time `json:"finished_at"`
	AccountID     int       `json:"account_id"`
}

type TransactionRepository interface {
	GetTransactionDetail(string) ([]Transaction, error)
	InsertTransaction(string, int) error
}

type TransactionHandler struct {
	Repo  TransactionRepository
	Store db.Store
}

func NewTransactionHandler(db db.Store) TransactionHandler {
	return TransactionHandler{Store: db}
}

func (transaction TransactionHandler) GetTransactionDetail(name string) ([]Transaction, error) {
	// defer transaction.Store.DB.Close()

	accId := 0

	rows, err := transaction.Store.DB.Query(fmt.Sprintf("SELECT account_id FROM account WHERE name = '%s' LIMIT 1", name))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&accId)
	}
	rows.Close()

	query := fmt.Sprintf("SELECT * FROM deposit_transaction WHERE account_id = %d", accId)
	rows2, err := transaction.Store.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var result []Transaction

	for rows2.Next() {
		var tempTransaction Transaction
		if err := rows2.Scan(&tempTransaction.TransactionID, &tempTransaction.Amount, &tempTransaction.Status, &tempTransaction.CreatedAt, &tempTransaction.FinishedAt, &tempTransaction.AccountID); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		result = append(result, tempTransaction)
	}

	return result, nil
}

func (transaction TransactionHandler) InsertTransaction(name string, amount int) error {
	// defer transaction.Store.DB.Close()
	accId := 0
	balance := 0

	rows, err1 := transaction.Store.DB.Query(fmt.Sprintf("SELECT account_id,balance FROM account WHERE name = '%s' LIMIT 1", name))
	if err1 != nil {
		return err1
	}

	for rows.Next() {
		rows.Scan(&accId, &balance)
	}
	rows.Close()

	query := fmt.Sprintf("INSERT INTO deposit_transaction (amount,status,created_at,finished_at,account_id) VALUES ('%d','success','%s','%s',%d)", amount, time.Now().String(), time.Now().String(), accId)
	_, err2 := transaction.Store.DB.Exec(query)
	if err2 != nil {
		return err2
	}

	update_query := fmt.Sprintf("UPDATE account SET balance=%d WHERE account_id = %d", balance+amount, accId)
	_, err3 := transaction.Store.DB.Exec(update_query)
	if err3 != nil {
		return err3
	}
	// defer rows.Close()
	return nil
}
