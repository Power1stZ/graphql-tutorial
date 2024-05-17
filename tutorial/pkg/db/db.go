package db

import (
	"database/sql"
	"fmt"
	"os"
)

type Method interface {
	// DBConnect() *sql.DB
}

type Store struct {
	DB     *sql.DB
	Method Method
}

func NewStore() Store {
	db, err := sql.Open("sqlite3", os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Println(err)
	}
	return Store{DB: db}
}

// func (s Store) GetData(q string,r *interface{}) (*interface{},error)
// func () DBConnect() *sql.DB {
// 	db, err := sql.Open("sqlite3", os.Getenv("DB_NAME"))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return db
// }
