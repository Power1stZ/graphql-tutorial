package jwt

import (
	"fmt"
	"graphql/tutorial/pkg/db"
)

type Method interface {
	CreateUser(string, string) error
}

type UserHandler struct {
	Method Method
	Store  db.Store
}

func NewUserHandler(db db.Store) UserHandler {
	return UserHandler{Store: db}
}

func (u UserHandler) CreateUser(name string, pass string) error {
	salt := GenerateRandomSalt(8)
	hash := HashString(name + pass + salt)

	query := `INSERT INTO user (user_name,password,salt,created_by) VALUES (?,?,?,?)`
	fmt.Println(name, pass)
	if _, err := u.Store.DB.Exec(query, name, hash, salt, "admin"); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u UserHandler) UserLogin(name string, pass string) bool {
	query := `SELECT * from user WHERE user_name = ? LIMIT 1`
	var user User
	err := u.Store.DB.QueryRow(query, name).Scan(&user.UserId, &user.Username, &user.Password, &user.Salt, &user.CreatedAt, &user.CreatedBy)
	if err != nil {
		fmt.Println(err)
		return false
	}

	checkPassword := HashString(name + pass + user.Salt)

	return user.Password == checkPassword
}
