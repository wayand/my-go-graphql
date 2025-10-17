# Go GraphQL with SQLite Example

A simple GraphQL API implementation in Go using SQLite as the database backend. This project demonstrates basic CRUD operations, GraphQL schema definition, queries, and mutations.

## 🚀 Features

- **GraphQL API** with queries and mutations
- **SQLite database** integration for data persistence
- **Book management system** with authors and descriptions
- **Type-safe GraphQL schema** definition
- Ready-to-run example queries and mutations

## 📋 Prerequisites

- Go 1.24.0 or later
- SQLite3 (for database inspection, optional)

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd my-go-graphql
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

## 📊 Database Schema

The application uses a simple SQLite database with the following schema:

```sql
CREATE TABLE books (
    id INTEGER,
    title TEXT
);
```

### Sample Data
The database comes pre-populated with:
- "Let's Go!" (ID: 1)
- "Let's Go Further!" (ID: 2) 
- "The Go Programming Language" (ID: 3)

## 🔧 GraphQL Schema

### Types

```graphql
type Book {
    id: Int
    title: String
    author: Author
    descriptions: [Description]
}

type Author {
    name: String
    books: [Int]
}

type Description {
    content: String
}
```

### Queries

- `book(id: Int)`: Fetch a specific book by ID from the database
- `list`: Fetch all books from the database

### Mutations

- `createBook(title: String!)`: Create a new book (in-memory only)

## 🚦 Usage Examples

The application demonstrates GraphQL operations with hardcoded examples:

### Query Example
```graphql
{
    book(id: 2) {
        id
        title
    }
}
```

### Mutation Example
```graphql
mutation {
    createBook(title: "Another Go book!") {
        title
    }
}
```

## 📁 Project Structure

```
.
├── main.go              # Main application with GraphQL setup
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── go-graphql-books.db  # SQLite database file
├── .gitignore           # Git ignore rules
└── README.md            # This file
```

## 🏗️ Architecture

### Data Layer
- **SQLite Database**: Stores book data persistently
- **In-Memory Data**: Used for complex objects with relationships (Author, Description)

### GraphQL Layer
- **Schema Definition**: Type definitions for Book, Author, Description
- **Resolvers**: Query and mutation handlers
- **Database Integration**: Direct SQL queries in resolvers

## 🔍 Key Components

### Structs
```go
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
```

### Database Operations
- Single book query with prepared statements
- List all books with result iteration
- Error handling for database operations

## 🧪 Testing

Run the application to see example output:

```bash
go run main.go
```

Expected output:
```
My Books, (a golang graphQL api)
{"data":{"createBook":{"title":"Another Go book!"}}}

{"data":{"book":{"id":2,"title":"Let's Go Further!"}}}
```

## 📦 Dependencies

- [`github.com/graphql-go/graphql`](https://github.com/graphql-go/graphql) - GraphQL implementation for Go
- [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) - SQLite3 driver for Go

## 🛣️ Potential Improvements

- [ ] Add HTTP server with GraphQL endpoint
- [ ] Implement proper error handling and validation
- [ ] Add authentication and authorization
- [ ] Extend database schema with foreign keys
- [ ] Add comprehensive test suite
- [ ] Implement pagination for list queries
- [ ] Add database migrations
- [ ] Docker containerization

## 📄 License

This is a learning/testing project. Feel free to use and modify as needed.

## 🤝 Contributing

This is an example project for testing purposes. Feel free to fork and experiment!