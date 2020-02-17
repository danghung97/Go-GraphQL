package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
)

type user struct{
	ID string `json:"id"`
	Name string `json:"name"`
}

var data map[string]user

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
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:         graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				idQuery, isOk := p.Args["id"].(string)
				if isOk {
					return data[idQuery], nil
				}
				return nil, nil
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var result []interface{}
				for _, value := range data {
					result = append(result, value)
				}
				return result, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	_ = importJSONDataFromFile("data.json", &data)
	
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={user(id:\"1\"){name}}'")
	http.ListenAndServe(":8080", nil)
}

func importJSONDataFromFile(fileName string, result *map[string]user) (isOk bool) {
	isOk = true
	content, err := ioutil.ReadFile(fileName)
	if err!=nil{
		fmt.Print("Error:", err)
		isOk = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOk = false
		fmt.Print("Error:", err)
	}
	return
}