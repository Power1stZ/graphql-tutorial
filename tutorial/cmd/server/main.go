package main

import (
	"fmt"
	"graphql/tutorial/internal/contact"
	"graphql/tutorial/pkg/graphqlfunc"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/mattn/go-sqlite3"
)

// Define GraphQL schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

// Define root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "Hello, GraphQL!", nil
			},
		},
		"contact": &graphql.Field{
			Type: graphql.NewList(graphqlfunc.ContactGraphQLType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return contact.GetContactData()
			},
		},
		"account": &graphql.Field{
			Type:    graphql.NewList(graphqlfunc.AccountGraphQLType),
			Args:    graphqlfunc.QueryArgs,
			Resolve: graphqlfunc.GetAccountResolver,
		},
		"transaction": &graphql.Field{
			Type:    graphql.NewList(graphqlfunc.DepositTransactionQLType),
			Args:    graphqlfunc.QueryArgs,
			Resolve: graphqlfunc.GetTransactionResolver,
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"insertAccount": &graphql.Field{
			Type:    graphql.NewList(graphqlfunc.AccountGraphQLType),
			Args:    graphqlfunc.InsertArgs,
			Resolve: graphqlfunc.InsertAccountResolver,
		},
		"insertTransaction": &graphql.Field{
			Type:    graphql.NewList(graphqlfunc.AccountGraphQLType),
			Args:    graphqlfunc.InsertArgs,
			Resolve: graphqlfunc.InsertTransactionResolver,
		},
	},
})

func main() {
	// Create a GraphQL handler for HTTP requests
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true,
	})

	// Serve GraphQL API at /graphql endpoint
	http.Handle("/graphql", graphqlHandler)

	// Start the HTTP server
	fmt.Println("Server is running at http://localhost:8000/graphql")
	http.ListenAndServe(":8000", nil)
}
