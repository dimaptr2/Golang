package acsengine

import (
	"os"
	"bufio"
	_ "github.com/nakagami/firebirdsql"
	"database/sql"
	"fmt"
	"time"
	"strings"
	"math/rand"
	"log"
	"strconv"
)

const confFile string = "resources/conf.txt"
const inputFile string = "resources/inputsql.txt"
const outputFile string = "resources/outputsql.txt"
const connection_string string = "SYSDBA:masterkey@srv-acs2.eatmeat.ru:3053/E:\\BASE\\OKO.FDB"

// Configuration data for the initial start
var slackers []string

// Panic at the disco
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// Read the initial configuration file
func ReadConf() error {

	file, err := os.Open(confFile)
	CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		slackers = append(slackers, scanner.Text())
	}

	return err
}

/*
	input - update input
	output - update output
	wait - it's clear from the context
 */

func defineActionType() string {

	// Here is the creation of the type of action.

	var action string

	datum := time.Now()
	hours, minutes, _:= datum.Clock()

	if (hours >= 8 && hours <= 17 && minutes > 0 && minutes <= 59) {
		action = "input"
	} else if (hours >= 18 && hours <= 19 && minutes > 0 && minutes <= 59) {
		action = "output"
	} else {
		action = "wait"
	}

	return action
}

func buildMonthMap() *map[string]string {

	m := make(map[string]string)

	m["Jan"] = "01"
	m["Feb"] = "02"
	m["Mar"] = "03"
	m["Apr"] = "04"
	m["May"] = "05"
	m["Jun"] = "06"
	m["Jul"] = "07"
	m["Aug"] = "08"
	m["Sep"] = "09"
	m["Oct"] = "10"
	m["Nov"] = "11"
	m["Dec"] = "12"
	// Give away the pointer to the map
	return &m
}

func getStringDate(pToTime *time.Time) string {

	// "Go" cannot read abstract symbols type of dd.MM.yyy HH:mm:ss
	// Only can read a concrete example.
	// Database has no ISO format of date.
	// Therefore, we convert ISO date to the internal format of the database.
	// And convert time type to the string type.

	dateStr := pToTime.Format(time.ANSIC)
	temp := strings.Split(dateStr, " ")
	dateStr = ""
	monthMap := buildMonthMap()
	for key, value := range *monthMap {
		if (temp[1] == key) {
			temp[1] = value
		}
	}

	// Insert current date into string
	dateStr = temp[2] + "." + temp[1] + "." + temp[4]

	return dateStr
}

// Build the SQL statement for the data selection
func buildSqlStatement(actType string) string {

	const space string = " "
	var fileName, sqlStmt string
	var sqlElements []string

	switch actType {
	case "input":
		fileName = inputFile
	case "output":
		fileName = outputFile
	case "wait":
		sqlStmt = "No"
	}

	if (sqlStmt != "No") {
		file, err := os.Open(fileName)
		CheckError(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			sqlElements = append(sqlElements, scanner.Text())
		}
		// Build the body of SQL command
		for _, value := range sqlElements {
			sqlStmt = sqlStmt + value + space
		}
	} //sqlStmt

	return sqlStmt
}

/* ------------------- Structure for analysis of data ---------------------------------------- */

// Declaration of the structure
type dbRowStructure struct {
	regId string
	userId string
	doorId string
	regDate string
	regDateFull string
}

var dbRow dbRowStructure
//var dbRows []dbRowStructure
var mapRows1, mapRows8 map[string][]dbRowStructure

// Create the map area for results of seeking
func makeInitialMap() {
	mapRows1 = make(map[string][]dbRowStructure)
	mapRows8 = make(map[string][]dbRowStructure)
}

// Set fields of structure (setter)
func (p *dbRowStructure) setDbRow(par *[]string) {

	if len(*par) == 5 {
		for index, value := range *par {
			switch index {
			case 0:
				p.regId = value
			case 1:
				p.userId = value
			case 2:
				p.doorId = value
			case 3:
				p.regDate = value
			case 4:
				p.regDateFull = value
			}
		}
	}
}

// Definition of the structure and realization of here methods

/* ------------------- End of structure definition ------------------------------------------- */

// Execute a query and fill the result set
func executeSqlQuery(db *sql.DB, command string) {

	var (
		regid, uid, doorid,
		regdate, regdatefull,
		dateAfter, dateBefore string
	)
        var counter int
	var flag bool

	// First of all, we get the current date into string format.
	datum := time.Now()
	currentDate := getStringDate(&datum)
	//layout := "18.05.2016, T15:04:05.000Z"
	dateAfter = currentDate + ", " + "00:01:00.000"
	dateBefore = currentDate + ", " + "23:59:59.000"
	makeInitialMap()

	for _, v := range slackers {
		tempSlice := []dbRowStructure{}
		resultSet, err := db.Query(command, v, dateAfter, dateBefore)
		CheckError(err)
		for resultSet.Next() {
			err := resultSet.Scan(&regid, &uid, &doorid,
						&regdate, &regdatefull)
			if err != nil {
				log.Fatal(err)
			}
			var params []string
			params = append(params, regid)
			params = append(params, uid)
			params = append(params, doorid)
			params = append(params, regdate)
			params = append(params, regdatefull)
			counter = len(params)
			dbRow.setDbRow(&params)
			tempSlice = append(tempSlice, dbRow)
			if doorid == "1" {
				flag = true
			} else {
				flag = false
			}
		}
		resultSet.Close()
		if counter == 5 {
			if (flag) {
				mapRows1[v] = tempSlice
			} else {
				mapRows8[v] = tempSlice
			} //if
		} //if
	} // external for

}

func displayResultSet(pm *map[string][]dbRowStructure) {

	const space  string = " "

	if len(*pm) > 0 {
		for key, _array := range *pm {
			fmt.Println(key + ":" + space)
			for _, rw := range _array {
				fmt.Print(rw.regId + space)
				fmt.Print(rw.userId + space)
				fmt.Print(rw.doorId + space)
				fmt.Print(rw.regDate + space)
				fmt.Print(rw.regDateFull + space)
				fmt.Print("\n")
			}
			fmt.Println("Number events: " + strconv.Itoa(len(_array)))
		}
	} else {
		fmt.Println("No data found")
	}

}

// Get random value
func getRandomValues(p_min, p_max int) int {
	return (rand.Intn(p_max - p_min) + p_min)
}

// Verify time intervals
func verifyAndChangeTime(pt *string, pact *string, door *string) {

	fmt.Println(*pt + "\t" + *pact + "\t" + *door)

	//layout := "2016-05-18T15:17:05:00Z"
	//datum, err := time.Parse(layout, *pt)
	//CheckError(err)
	//
	//year, _, day := datum.Date()
	//monthStr := datum.Month().String()
	//// Build the month as string
	//switch monthStr {
	//case "January":
	//	monthStr = "01"
	//case "February":
	//	monthStr = "02"
	//case "March":
	//	monthStr = "03"
	//case "April":
	//	monthStr = "04"
	//case "May":
	//	monthStr = "05"
	//case "June":
	//	monthStr = "06"
	//case "July":
	//	monthStr = "07"
	//case "August":
	//	monthStr = "08"
	//case "September":
	//	monthStr = "09"
	//case "October":
	//	monthStr = "10"
	//case "November":
	//	monthStr = "11"
	//case "December":
	//	monthStr = "12"
	//}
	//hours, minutes, seconds := datum.Clock()
	//
	//*pt = strconv.Itoa(day) + "." + monthStr +
	//"." + strconv.Itoa(year)
	//*pt = *pt + ", "
	//
	//switch (*pact) {
	//case "input":
	//	if hours >= 8 && minutes >= 35 &&
	//		seconds >= 0 && seconds <= 59 {
	//		hours = 8
	//		minutes = getRandomValues(25, 35)
	//		if *door == "8" {
	//			seconds += 10
	//		}
	//	}
	//case "output":
	//	if hours <= 17 && minutes < 30 &&
	//	seconds >= 0 && seconds <= 59 {
	//		hours = 17
	//		minutes = getRandomValues(30, 40)
	//		if *door == "8" {
	//			seconds -= 10
	//		}
	//	}
	//}
	//
	//*pt = *pt + strconv.Itoa(hours) + ":" + strconv.Itoa(minutes) +
	//":" + strconv.Itoa(seconds) + ".000"

}

// Make the changing
func changeDbState(db *sql.DB, pm *map[string][]dbRowStructure,
			pact string) error {

	var _err error

	// Start transaction
	//tx, err := db.Begin()
	//_err = err
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer tx.Rollback()

	if len(*pm) > 0 {
		//stmt, err := db.Prepare("UPDATE graph_fact_events " +
		//"SET regdate = ?, regdatefull = ? " +
		//"WHERE regid = ? AND uid = ?")
		//_err = err
		//defer stmt.Close()
		// events
		for key, _array := range *pm {
			if len(_array) > 0 {
				for _index, _value := range _array {
					if _index == 0 || _index == (len(_array) - 1) {
						if key == _value.userId {
							//fmt.Println(_value.regDate)
							verifyAndChangeTime(&_value.regDate,
								&pact, &_value.doorId)
							_value.regDateFull = _value.regDate
							//_, err := stmt.Exec(_value.regDate, _value.regDateFull,
							//	_value.regId, _value.userId)
							//if err != nil {
							//	_err = err
							//	log.Fatal(err)
							//}
						}
					}
				} // for
			}

		} //for
		// total intervals
		//_err = tx.Commit()
	}

	return _err
}

// clearing the garbage
func verifyDataForUpdating(pm *map[string][]dbRowStructure) {

	for _, _array := range *pm {
		pta := &_array
		for i := 0; i < len(_array); i++ {
			*pta = append(_array[:i], _array[i+1:]...)
		}
	}

}

func makeChanges(db *sql.DB, pm *map[string][]dbRowStructure,
			pact string) error {

	var _err error

	if len(*pm) > 0 {
		verifyDataForUpdating(pm)
		_err = changeDbState(db, pm, pact)
		CheckError(_err)
	} else {
		log.Println("No data for changing")
	}

	return _err
}

// Central function, that can process data.
/*	Parameter "mode":
	R - only read (data selection)
  	U - updating
*/
func ProcessDbTask(mode string) error {

	var _err error

	dbConn, err := sql.Open("firebirdsql", connection_string )
	_err = err
	CheckError(err)
	defer dbConn.Close()

	fmt.Println("That's all right!")
	act := defineActionType()
	sqlFind := buildSqlStatement(act)
	if act == "wait" {
		return _err
	}

	executeSqlQuery(dbConn, sqlFind)

	switch mode {
	case "R":
		fmt.Println("Door Id: 1")
		displayResultSet(&mapRows1)
		fmt.Println("Door Id: 8")
		displayResultSet(&mapRows8)
	case "U":
		// doorId == 1
		_err = makeChanges(dbConn, &mapRows1, act)
		// doorId == 8
		_err = makeChanges(dbConn, &mapRows8, act)
	}

	return _err
}
