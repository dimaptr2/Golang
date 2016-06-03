package bengi

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//	"strconv"
	"strings"
	// SAP connection through RFC
	//	"github.com/sap/gorfc/gorfc"

	"database/sql"
	// SQlite3 Go API
	_ "github.com/mattn/go-sqlite3"
	//"time"
	"strconv"
)

// Types

// type of rows for the database updating

type rowDbType struct {
	year      string
	month     string
	day       string
	dayweek   string
	eventtime string
}

// Global parameters of this package
const fPath string = "./beaglebone"
const c_dbname string = "beagledb.db"

var catalog []os.FileInfo
var contentFile []string

//var sapconn gorfc.Connection
var result [][]string

// main functions

// Get list of files in the directory
func getCatalogList(fpath string) []os.FileInfo {

	files, err := ioutil.ReadDir(fpath)
	if err != nil {
		log.Fatal(err)
	}
	return files

}

// read the file line by line
func readEventFiles(files []os.FileInfo) {

	// Read the catalog
	for _, file := range files {

		f, err := os.Open(fPath + "/" + file.Name())
		if err != nil {
			panic(err)
		}


		s := bufio.NewScanner(f)
		//		 Read the file
		for s.Scan() {
			contentFile = append(contentFile, s.Text())
		} //second cycle
		f.Close()

	} // first cycle

} //readEventFiles

func readContentOfFile(pcache *[][]string) {

	var numberElems int = len(contentFile)

	fmt.Print("Number rows: ")
	fmt.Println(numberElems)

	for i := 0; i < numberElems; i++ {
		temp := strings.Split(contentFile[i], " ")
		*pcache = append(*pcache, temp)
	}

	fmt.Println("Resulting array contains rows: " + strconv.Itoa(len(*pcache)))

}

// Create the table in the database if it doesn't exist
func createTable(db *sql.DB) error {

	xsql := `
	create table events (year int not null, month varchar(3) not null,
	day int not null, dayweek varchar(3) not null,
	eventtime time not null,
	primary key(year, month, day, dayweek, eventtime));`
	_, err := db.Exec(xsql)
	if err != nil {
		log.Printf("%q: %s\n", err, xsql)
	}
	return err
}

// Create the database file
func createDB() error {

	sqdb, err := sql.Open("sqlite3", c_dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer sqdb.Close()

	err = createTable(sqdb)
	if err != nil {
		log.Fatal(err)
	}
	return err

}

// SQLite database processing ...
func processDb(pdata *[][]string) {

	var line rowDbType
	var ptr1 *[]string
	var flag bool
	var counter int = 0
	//const cpath string = "/opt/cachebeagle/"

	pdb, err := sql.Open("sqlite3", c_dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer pdb.Close()

	// Build sql
	sqlStmt := "INSERT INTO events VALUES (?, ?, ?, ?, ?)"

	// start transaction
	tx, err := pdb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	// Run through collection of arrays and create the statements
	for _, v1 := range *pdata {
		ptr1 = &v1
		// If the length of row is not equal 12 bytes (6 symbols of UTF-8),
		// then go to the next iteration.
		flag = false
			for index2, v2 := range *ptr1 {

					switch index2 {
					case 1:
						line.dayweek = v2
					case 2:
						line.month = v2
					case 3:
						line.day = v2
					case 4:
						line.eventtime = v2
					case 5:
						line.year = v2
					}
			} // second for
		if line.year != "" {
			flag = true
		}
		//// database updating
		if flag {
			_, err = stmt.Exec(line.year, line.month, line.day,
				line.dayweek, line.eventtime)
			if err != nil {
				//log.Fatal(err)
				continue
			} else {
				counter++
			}
		}
	} // first for

	tx.Commit() // end of data loading and commit work
	fmt.Println("Success! Database was processed...")
	fmt.Printf("Into database %d rows was inserted...\n", counter)

}

// Central function that can run the task
func RunTaskBeagleSrv() {

	catalog = getCatalogList(fPath)
	readEventFiles(catalog)
	readContentOfFile(&result)
	_, err := os.Stat(c_dbname)
	if err != nil {
		err = createDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	processDb(&result)

}
