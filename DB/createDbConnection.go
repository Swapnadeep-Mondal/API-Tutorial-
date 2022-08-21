package conn

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // creating a global variable so that it can be accessed in other package/ function

func Creating_connection() { // function name should start with capital letter so that it becomes visible outside the package

	fmt.Println("Starting of the function")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		"root", "secretpassword", "127.0.0.1:3306", "testdb",
	)
	fmt.Println("dataSourceName", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err.Error())
	}

	// defer db.Close()

	DB = db
	if err = DB.Ping(); err != nil { // to check connection is working or not
		panic(err)
	}

	fmt.Println("Sucessful Connection to Database")
}
