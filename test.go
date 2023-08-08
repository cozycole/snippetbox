package main

import (
	"fmt"
	"path/filepath"

	"log"
	// "github.com/go-sql-driver/mysql"
	// "database/sql"
)

func main() {
	// fmt.Print(path.Clean("/../har_parsing") + "\n")
	fmt.Print(log.Ltime)
	fpaths, _ := filepath.Glob("./ui/html/*/*")
	for x, p := range fpaths {
		fmt.Print(x, p, "\n")
	}
	// db, err := sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")
	// if err != nil {
	// 	return
	// }
	// db = nil
}
