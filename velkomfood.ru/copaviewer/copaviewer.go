package main

import (
	"fmt"
	"velkomfood.ru/copaviewer/copavwlib"
)

// Start point of this server
func main() {

	fmt.Println("Start COPA viewer")

	_err := copavwlib.OpenSAPConnection()
	if (_err != nil) {
		fmt.Println("SAP connection is failed!")
		panic(_err)
	} else {
		copavwlib.CloseSAPConnection()
	}


}
