package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID           int
	Title        string
	Author       Author
	Descriptions []Description
}

type Author struct {
	Name  string
	Books []int
}

type Description struct {
	Content string
}

func populate() []Book {
	author := &Author{Name: "Alex Edwards", Books: []int{1, 2}}
	book := Book{
		ID:     1,
		Title:  "Let's Go!",
		Author: *author,
		Descriptions: []Description{
			Description{Content: "it's a good book."},
		},
	}
	book2 := Book{
		ID:     2,
		Title:  "Let's Go Further!",
		Author: *author,
		Descriptions: []Description{
			Description{Content: "it's a good book too."},
		},
	}
	book3 := Book{
		ID:     3,
		Title:  "Let's Go Further2!",
		Author: *author,
		Descriptions: []Description{
			Description{Content: "it's a good book too2."},
		},
	}
	var books []Book
	books = append(books, book)
	books = append(books, book2)
	books = append(books, book3)

	return books
}

func main() {
	fmt.Println("My Books, (a golang graphQL api)")
	books := populate()

	var descriptionType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Description",
		Fields: graphql.Fields{
			"content": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var authorType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"books": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	var bookType = graphql.NewObject(graphql.ObjectConfig{
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
			"descriptions": &graphql.Field{
				Type: graphql.NewList(descriptionType),
			},
		},
	})

	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createBook": &graphql.Field{
				Type:        bookType,
				Description: "Create a new book",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					book := Book{
						ID:    len(books) + 1,
						Title: p.Args["title"].(string),
					}
					books = append(books, book)
					return book, nil
				},
			},
		},
	})

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
					db, err := sql.Open("sqlite3", "./go-graphql-books.db")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()
					var book Book
					err = db.QueryRow("SELECT ID, Title FROM books where ID = ?", id).Scan(&book.ID, &book.Title)
					if err != nil {
						fmt.Println(err)
					}
					return book, nil
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(bookType),
			Description: "Get list of All Books.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, err := sql.Open("sqlite3", "./go-graphql-books.db")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				var books []Book
				results, err := db.Query("SELECT * FROM books")
				if err != nil {
					fmt.Println(err)
				}

				for results.Next() {
					var book Book
					err = results.Scan(&book.ID, &book.Title)
					if err != nil {
						fmt.Println(err)
					}
					log.Println(book)
					books = append(books, book)
				}

				return books, nil
			},
		},
	}

	// defines the object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// defines a chema config
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	// creates out schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL schema, err %v", err)
	}

	mutationQuery := `
		mutation {
			createBook(title: "Another Go book!") {
				title
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: mutationQuery}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)

	query := `
		{
			book (id: 2) {
				id
				title
			}
		}
	`

	params = graphql.Params{Schema: schema, RequestString: query}
	r = graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ = json.Marshal(r)
	fmt.Println("\n")
	fmt.Printf("%s \n", rJSON)
}
