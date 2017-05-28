package copavwlib

import (
	"database/sql"
	"simonwaldherr.de/go/saprfc"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
)

// Basic variables
var sapConnection *saprfc.Connection
var db *sql.DB

// Main data structures
type OrganizationalStructure struct {
	CompanyCode string
	ControllingArea string
	Plant string
	SalesOrganization string
}

type Channel struct {
	Id string
	Description string
}

type Sector struct {
	Id string
	Description string
}

type Region struct {
	Id string
	Description string
}

type Material struct {
	Id uint32
	Description string
}

type Customer struct {
	Id uint32
	Name string
}

type Account struct {
	Id uint32
	Description string
}

type CostCenter struct {
	Id uint32
	Description string
}

type InternalOrder struct {
	Id string
	Description string
}

type CE1 struct {
	Belnr uint32
	Posnr uint32
	Operation string
	Period string
	Hzdat time.Time
	Budat time.Time
	Currency string
	UnitOfMeasure string
	ChannelId string
	SectorId string
	RegionId string
	CustomerId uint32
	MaterialId uint32
	InvoiceType string
	InvoiceId uint32
	AccountId uint32
	Quantity decimal.Decimal
	Turnover decimal.Decimal
}

// Basic functions and methods

func abapSystem() saprfc.ConnectionParameter {
	return saprfc.ConnectionParameter{
		Dest:      "XXX",
		Client:    "XX",
		User:      "XXXXXX",
		Passwd:    "XXXXXXX",
		Lang:      "XX",
		Ashost:    "xx.xx.xx.xx",
		Sysnr:     "xx",
		Sysid:     "XXX",
		Saprouter: "XXXXXXXXX",
	}

}
// Open and close a connection to the SAP server
func OpenSAPConnection() error {
	var err error
	sapConnection, err = saprfc.ConnectionFromParams(abapSystem())
	return err
}

func CloseSAPConnection() {
	sapConnection.Close()
}

// Open and close a connection to the SQLite database file
func OpenDBConnection() {

}

func CloseDBConnection() {
	db.Close()
}


