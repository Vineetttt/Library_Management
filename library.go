package main

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// struct that stores library information
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// slice named books of type book (struct)
var books = []book{
	{ID: "1", Title: "Think Like a Monk", Author: "Jay Shetty", Quantity: 2},
	{ID: "2", Title: "The Power of your subconscious mind", Author: "Joseph Murphy", Quantity: 5},
	{ID: "3", Title: "Think and Grow Rich", Author: "Napoleon Hill", Quantity: 6},
}

// function to return a JSON version of the books
func getBooks(c *gin.Context) {
	c.IndentedJSON(200, books)
}

// function to update the library (add books)
func addBook(c *gin.Context) {
	var newBook book

	// if we encounter any error while binding our newBook info with the current books slice
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	// append newBook info to the books slice
	books = append(books, newBook)
	c.IndentedJSON(201, newBook)
}

// helper function to get a book by id
// takes id as input and returns the book and error (if)
func getBookById(id string) (book, error) {
	for i, b := range books {
		if b.ID == id {
			return books[i], nil
		}
	}
	return book{}, errors.New("book does not exist")
}

// function to get book by id
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	// if the book does not exist
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(200, book)
}

// function to decrease the quantity of book by 1
func decreaseBookQty(c *gin.Context) {
	id, success := c.GetQuery("id")

	// if id is not passed in the route i.e. success == false
	if !success {
		c.IndentedJSON(400, gin.H{"message": "book id parameter missing"})
		return
	}
	// if the id is found we try to search for the book in the library
	// if the id is found we try to search for the book in the library
	for i := range books {
		if books[i].ID == id {
			// if the book exists but is not available (quantity <= 0)
			if books[i].Quantity <= 0 {
				c.IndentedJSON(400, gin.H{"message": "book not available"})
				return
			}

			// otherwise, decrease the quantity of the book by 1
			books[i].Quantity -= 1
			c.IndentedJSON(200, books[i])
			return
		}
	}

	// if the book does not exist
	c.IndentedJSON(404, gin.H{"message": "book not found"})
}

// function to increase the book quantity by 1
func increaseBookQty(c *gin.Context) {
	id, success := c.GetQuery("id")

	// if id is not passed in the route i.e. success == false
	if !success {
		c.IndentedJSON(400, gin.H{"message": "book id parameter missing"})
		return
	}

	// if the id is found, search for the book in the library
	for i := range books {
		if books[i].ID == id {
			// increase the quantity of the book by 1
			books[i].Quantity += 1
			c.IndentedJSON(200, books[i])
			return
		}
	}

	// if the book does not exist
	c.IndentedJSON(404, gin.H{"message": "book not found"})
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", addBook)
	router.PATCH("/checkout", decreaseBookQty)
	router.PATCH("/return", increaseBookQty)
	router.Run(":8080")
}
