package graphqlfunc

import (
	"fmt"
	"graphql/tutorial/internal/deposit"
	"graphql/tutorial/internal/jwt"
	"graphql/tutorial/pkg/db"
	"log"

	"github.com/graphql-go/graphql"
)

func GetAccountResolver(p graphql.ResolveParams) (interface{}, error) {
	store := db.NewStore()
	account := deposit.NewAccountHandler(store)

	if _, ok := p.Args["name"].(string); !ok {
		log.Fatal("Input not Valid")
	}
	return account.GetAccountDetail(p.Args["name"].(string))
}

func GetTransactionResolver(p graphql.ResolveParams) (interface{}, error) {
	store := db.NewStore()
	transaction := deposit.NewTransactionHandler(store)

	if _, ok := p.Args["name"].(string); !ok {
		log.Fatal("Input not Valid")
	}
	return transaction.GetTransactionDetail(p.Args["name"].(string))
}

func InsertAccountResolver(p graphql.ResolveParams) (interface{}, error) {
	name := ""
	if v, ok := p.Args["name"].(string); ok {
		name = v
	}
	store := db.NewStore()
	account := deposit.NewAccountHandler(store)
	if err := account.InsertAccount(name, 0); err != nil {
		fmt.Println(err.Error())
	}
	return account.GetAccountDetail(name)
}

func InsertTransactionResolver(p graphql.ResolveParams) (interface{}, error) {
	amount := 0
	if v, ok := p.Args["amount"].(int); ok {
		amount = v
	}
	name := ""
	if v, ok := p.Args["name"].(string); ok {
		name = v
	}
	store := db.NewStore()
	transaction := deposit.NewTransactionHandler(store)
	if err := transaction.InsertTransaction(name, amount); err != nil {
		fmt.Println(transaction)

		fmt.Println(err.Error())
	}
	return transaction.GetTransactionDetail(name)
}

func CreateUserResolver(p graphql.ResolveParams) (interface{}, error) {
	store := db.NewStore()
	user := jwt.NewUserHandler(store)
	name := ""
	if v, ok := p.Args["user_name"].(string); ok {
		name = v
	}
	pass := ""
	if v, ok := p.Args["password"].(string); ok {
		pass = v
	}

	if error := user.CreateUser(name, pass); error != nil {
		return "fail", error
	}
	return "success", nil
}

func UserLoginResolver(p graphql.ResolveParams) (interface{}, error) {
	store := db.NewStore()
	user := jwt.NewUserHandler(store)
	name := ""
	if v, ok := p.Args["user_name"].(string); ok {
		name = v
	}
	pass := ""
	if v, ok := p.Args["password"].(string); ok {
		pass = v
	}
	ok := user.UserLogin(name, pass)
	return ok, nil
}
