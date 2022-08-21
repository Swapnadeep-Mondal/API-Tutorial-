package main

import (
	conn "Go-Api-First/DB"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type badRequest struct {
	Msg string `json:"message"`
	Err string `json:"error"`
}

var books = []book{
	{ID: 1, Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: 3, Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

var final_book []book

func getBooks(c *gin.Context) {

	fmt.Println("Entering getBooks function")
	result, err := conn.DB.Query("select * from books limit 5")

	if err != nil {
		var err1 badRequest
		err1.Err = err.Error()
		err1.Msg = ""
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	for result.Next() {
		var tag book
		err1 := result.Scan(&tag.ID, &tag.Title, &tag.Author, &tag.Quantity)

		if err1 != nil {
			var err2 badRequest
			err2.Err = err1.Error()
			err2.Msg = "Error occured in retrieving data"
			c.JSON(http.StatusUnauthorized, err2)
			return
		}

		final_book = append(final_book, tag) // to store every query row in final_book

		/*
			log.Println(tag.ID)
			log.Printf(tag.Title)
			log.Println(tag.Author)
			log.Println(tag.Quantity)
		*/
	}

	c.IndentedJSON(http.StatusOK, final_book)
	// c.IndentedJSON(http.StatusOK, books)
	return
}

func getBookDetailsById(c *gin.Context) { // localhost:8080/bookbyid/1

	/*
		fmt.Println(c.Param("book_id"))
		id_book := c.Param("book_id")
		for _, b := range books {
			id, err := strconv.Atoi(id_book)

			if err != nil {
				c.JSON(http.StatusBadRequest, "Invalid Book ID")
				return
			}

			if b.ID == id {
				fmt.Println(b)
				c.JSON(http.StatusOK, b)
				return
			}
		}

		c.JSON(http.StatusOK, "Book ID not found")
		return
	*/

	id_book := c.Param("book_id")
	var tag book

	err := conn.DB.QueryRow("select * from books where ID=?", id_book).Scan(&tag.ID, &tag.Title, &tag.Author, &tag.Quantity)

	if err != nil {
		var err1 badRequest
		err1.Err = err.Error()
		err1.Msg = "Data not found in the Table "
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	c.IndentedJSON(http.StatusOK, tag)
	return
}

func getBookDetailsByEnteringId(c *gin.Context) { // GET ---->  localhost:8080/bookbyid/id=1

	/*
		id_book := c.Param("enter_book_id")
		for _, b := range books {
			id, err := strconv.Atoi(id_book)

			if err != nil {
				c.JSON(http.StatusBadRequest, "Invalid Book ID")
				return
			}

			if b.ID == id {
				fmt.Println(b)
				c.JSON(http.StatusOK, b)
				return
			}
		}

		c.JSON(http.StatusOK, "Book ID is invalid")
		return
	*/

	id_book := c.Param("enter_book_id")
	var tag book

	err := conn.DB.QueryRow("select * from books where ID=?", id_book).Scan(&tag.ID, &tag.Title, &tag.Author, &tag.Quantity)

	if err != nil {
		var err1 badRequest
		err1.Err = err.Error()
		err1.Msg = "Table does not have any row with this ID"
		c.IndentedJSON(http.StatusBadRequest, err1)
		return
	}

	c.IndentedJSON(http.StatusOK, tag)
	return
}

func getBookDetailsThroughBody(c *gin.Context) {
	var viewBook book

	if err := c.BindJSON(&viewBook); err != nil { // when checking for id=1,2,3 , err==nil. So checking the condition nil!=nil, which is false, it doesnt enter the loop
		//	fmt.Println("ERRORRRRRRRR!!!!")
		// 	fmt.Println("The error is ", err)

		var err1 badRequest
		err1.Err = err.Error()
		err1.Msg = "Invalid Request"
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	for _, b := range books {
		if b.ID == viewBook.ID {
			c.JSON(http.StatusOK, b)
			return
		}
	}

	var err1 badRequest
	err1.Err = ""
	err1.Msg = "Not found"
	c.JSON(http.StatusNotFound, err1)
	return

}

func createBook(c *gin.Context) { // POST ----> 	localhost:8080/bookss
	var newBook book
	fmt.Println(newBook)

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"mesaage": "", "error": err.Error()})
		return
	}
	fmt.Println(newBook)

	result, err := conn.DB.Exec("insert into books values (?,?,?,?)", newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)

	if err != nil {
		var err1 badRequest
		err1.Err = err.Error()
		err1.Msg = "Invalid Data Entry Attempted in Database"
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	fmt.Println(result.LastInsertId())

	// books = append(books, newBook)
	c.IndentedJSON(http.StatusOK, newBook)
}

func main() {

	conn.Creating_connection()

	router := gin.Default()
	router.GET("/books", getBooks)

	router.GET("/bookbyid/:book_id", getBookDetailsById)

	router.GET("/bookbyid/id=:enter_book_id", getBookDetailsByEnteringId)

	router.GET("/booksbody", getBookDetailsThroughBody)

	router.POST("/bookss", createBook)
	router.Run("localhost:8080")
}
