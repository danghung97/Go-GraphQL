package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/labstack/gommon/log"
	"net/http"
)

var Schema graphql.Schema

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Fields:      graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Query",
	Fields:      graphql.Fields{
		"me": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Context.Value("currentUser"), nil
			},
		},
	},
})

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	user := struct{
		ID int `json:"id"`
		Name string `json:"name"`
	}{1, "Dang Hung"}
	
	result := graphql.Do(graphql.Params{
		Schema:         Schema,
		RequestString:  r.URL.Query().Get("query"),
		Context:        context.WithValue(context.Background(), "currentUser", user),
	})
	
	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={me{id,name}}'")
	err := http.ListenAndServe(":8080", nil)
	if err!=nil {
		log.Fatal(err)
		return
	}
	
}

func init() {
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	
	if err!=nil {
		log.Fatalf("failed to create schema, error: %v", err)
	}
	
	Schema = s
}