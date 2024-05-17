package deposit

import (
	"fmt"
	"graphql/tutorial/pkg/db"
	"log/slog"
)

type Account struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Balance   int    `json:"balance"`
}

type AccountRepository interface {
	GetAccountDetail(string) ([]Account, error)
	InsertAccount(string, int) error
}

type AccountHandler struct {
	Repo  AccountRepository
	Store db.Store
}

func NewAccountHandler(db db.Store) AccountHandler {
	return AccountHandler{Store: db}
}

func (account AccountHandler) GetAccountDetail(name string) ([]Account, error) {
	// db, err := db.DbConnection()
	// if err != nil {
	// 	return nil, err
	// }
	// defer account.Store.DB.Close()

	query := fmt.Sprintf("SELECT * FROM account WHERE name = '%s'", name)
	rows, err := account.Store.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Account

	for rows.Next() {
		var tempAccount Account
		if err := rows.Scan(&tempAccount.AccountID, &tempAccount.Name, &tempAccount.Status, &tempAccount.Balance); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		result = append(result, tempAccount)
	}

	return result, nil
}

func (account AccountHandler) InsertAccount(name string, amount int) error {

	// defer account.Store.DB.Close()

	query := fmt.Sprintf("INSERT INTO account (name,status,balance) VALUES ('%s','active',%d)", name, amount)
	_, err := account.Store.DB.Exec(query)
	if err != nil {
		return err
	}
	// defer rows.Close()
	return nil
}
