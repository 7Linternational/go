// 1. Connect with postgresql DONE
// 2. Fetch data DONE
// 3. Update data DONE

package main

import (
	// "fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

/**
 * [checkErr description]
 * @return {[type]} [description]
 */
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * [checkErrCallback description]
 * @param  {[type]} err error         [description]
 * @param  {[type]} fn  func()        [description]
 * @return {[type]}     [description]
 */
func checkErrCallback(err error, fn func()) {
	if err != nil {
		if fn != nil {
			fn()
		}
		log.Fatal(err)
	}
}

/**
 * [getData description]
 * @return {[type]} [description]
 */
func getData(db *sql.DB) {
	var (
		id      int
		name    string
		address string
	)

	rows, err := db.Query("SELECT id, name, address FROM company;")

	checkErr(err)

	// Release db connecton to pool
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &address)

		checkErr(err)

		log.Println(id, name, address)
	}
}

/**
 * [modifyData description]
 * @param  {[type]} db *sql.DB       [description]
 * @return {[type]}    [description]
 */
func modifyData(db *sql.DB) {
	var stmt *sql.Stmt
	var res sql.Result

	stmt, err := db.Prepare("INSERT INTO company(name, address) VALUES ($1,$2)")
	checkErr(err)

	res, err = stmt.Exec("Carnevale Interactive", "Grand Rapids MI")
	checkErr(err)

	lastId, err := res.LastInsertId()
	checkErr(err)
	log.Println(lastId)
}

/**
 * [modifyDataTrans description]
 * @param  {[type]} db *sql.DB       [description]
 * @return {[type]}    [description]
 */
func modifyDataTrans(db *sql.DB) {
	Tx, err := db.Begin()
	checkErr(err)
	defer Tx.Rollback()

	var stmt *sql.Stmt

	stmt, err = Tx.Prepare("INSERT INTO company(name, address) VALUES ($1,$2)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec("Carnevale Interactive", "Grand Rapids MI")
	checkErr(err)

	err = Tx.Commit()
	checkErr(err)
}

/**
 * [main description]
 * @return {[type]} [description]
 */
func main() {
	db, err := sql.Open("postgres", "dbname=speclint user=postgres password=japan4 port=5432 sslmode=disable")
	checkErrCallback(err, func() {
		log.Fatal("here")
	})

	// Test db connection
	err = db.Ping()
	checkErr(err)

	getData(db)

	// modifyData(db)

	modifyDataTrans(db)

	// Close connection
	defer db.Close()
}
