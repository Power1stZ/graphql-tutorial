package graphqlfunc

import "github.com/graphql-go/graphql"

var QueryArgs = graphql.FieldConfigArgument{
	"name": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var InsertArgs = graphql.FieldConfigArgument{
	"name": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"amount": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var UserArgs = graphql.FieldConfigArgument{
	"user_name": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"password": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}
