package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Book struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name  string
	Books []int
}

type Comment struct {
	Content string
}

func populate() []Book {
	author := &Author{Name: "Alex Edwards", Books: []int{1,2}}
	book := Book{
		ID:     1,
		Title:  "Let's Go!",
		Author: *author,
		Comments: []Comment{
			Comment{Content: "it's a good book."},
		},
	}
	book2 := Book{
		ID:     2,
		Title:  "Let's Go Further!",
		Author: *author,
		Comments: []Comment{
			Comment{Content: "it's a good book too."},
		},
	}
	var books []Book
	books = append(books, book)
	books = append(books, book2)

	return books
}

func main() {
	fmt.Println("My Books, (a golang graphQL api)")
	books := populate()

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"Content": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Books": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var bookType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Book",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"book": &graphql.Field{
			Type:        bookType,
			Description: "Get Book by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, book := range books {
						if int(book.ID) == id {
							return book, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(bookType),
			Description: "Get list of All Books.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return books, nil
			},
		},
	}

	// defines the object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// defines a chema config
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// creates out schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL schema, err %v", err)
	}

	query := `
		{
			book(id:2) {
				title
				author {
					Name
					Books
				}
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
