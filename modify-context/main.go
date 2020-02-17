package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
)

type User struct {
	ID int `json:"id"`
}

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Fields:      graphql.Fields{
		"id": &graphql.Field{
			Type:              graphql.Int,
			Resolve:           func(p graphql.ResolveParams) (interface{}, error){
				rootvalue := p.Info.RootValue.(map[string]interface{})
				if rootvalue["data-from-parent"] == "ok" &&
					rootvalue["data-before-execution"] == "ok" {
						user := p.Source.(User)
						return user.ID, nil
				}
				return nil, nil
			},
		},
	},
})

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Query",
			Fields:      graphql.Fields{
				"users": &graphql.Field{
					Type:              graphql.NewList(UserType),
					Resolve:           func(p graphql.ResolveParams) (interface{}, error) {
						rootvalue := p.Info.RootValue.(map[string]interface{})
						rootvalue["data-from-parent"] = "ok"
						result := []User{
							{ID: 2},
							{ID: 3},
						}
						return result, nil
					},
				},
			},
		}),
	})
	rootObject := map[string]interface{}{
		"data-before-execution": "ok",
	}
	
	result := graphql.Do(graphql.Params{
		Schema: schema,
		RootObject: rootObject,
		RequestString: "{ users { id } }",
	})
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
}