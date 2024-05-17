package graphqlfunc

import (
	"github.com/graphql-go/graphql"
)

var ContactGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"contact_id":   &graphql.Field{Type: graphql.Int},
		"name":         &graphql.Field{Type: graphql.String},
		"first_name":   &graphql.Field{Type: graphql.String},
		"last_name":    &graphql.Field{Type: graphql.String},
		"gender_id":    &graphql.Field{Type: graphql.Int},
		"dob":          &graphql.Field{Type: graphql.String},
		"email":        &graphql.Field{Type: graphql.String},
		"phone":        &graphql.Field{Type: graphql.String},
		"address":      &graphql.Field{Type: graphql.String},
		"photo_path":   &graphql.Field{Type: graphql.String},
		"created_date": &graphql.Field{Type: graphql.DateTime},
		"created_by":   &graphql.Field{Type: graphql.String},
	},
})

var AccountGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"account_id": &graphql.Field{Type: graphql.Int},
		"name":       &graphql.Field{Type: graphql.String},
		"status":     &graphql.Field{Type: graphql.String},
		"balance":    &graphql.Field{Type: graphql.Int},
	},
})

var DepositTransactionQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DepositTransaction",
	Fields: graphql.Fields{
		"trans_id":    &graphql.Field{Type: graphql.Int},
		"amount":      &graphql.Field{Type: graphql.Int},
		"status":      &graphql.Field{Type: graphql.String},
		"created_at":  &graphql.Field{Type: graphql.DateTime},
		"finished_at": &graphql.Field{Type: graphql.DateTime},
		"account_id":  &graphql.Field{Type: graphql.Int},
	},
})

var UserQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserLogin",
	Fields: graphql.Fields{
		"user_name":  &graphql.Field{Type: graphql.String},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"created_by": &graphql.Field{Type: graphql.String},
	},
})
