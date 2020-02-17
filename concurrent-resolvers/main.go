package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

type Foo struct {
	Name string
}

var FieldFooType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Food",
	Fields:      graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

type Bar struct {
	Name string
}

var FieldBarType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Bar",
	Fields:      graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"concurrentFieldFoo": &graphql.Field{
			Type: FieldFooType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var foo = Foo{Name: "Foo's Name"}
				return func() (interface{}, error) {
					return &foo, nil
				}, nil
			},
		},
		"concurrentFieldBar": &graphql.Field{
			Type: FieldBarType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var bar = Bar{Name: "Bar's Name"}
				return func() (interface{}, error) {
					return &bar, nil
				}, nil
			},
		},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryType,
	})
	if err!=nil {
		log.Fatal(err)
	}
	query := `
		query {
			concurrentFieldFoo {
				name
			}
			concurrentFieldBar {
				name
			}
		}
	`
	
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema: schema,
	})
	b, err := json.Marshal(result)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
}
