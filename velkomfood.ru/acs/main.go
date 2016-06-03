package main

import (
	"velkomfood.ru/acs/acsengine"
	"fmt"
	"os"
	"log"
)

func main() {

	err := acsengine.ReadConf()
	if err != nil {
		fmt.Println("Cannot open the file conf.txt")
		os.Exit(1)
	} else {
		if len(os.Args) == 2 {
			err := acsengine.ProcessDbTask(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Enter R (Read) or U (Update")
			os.Exit(1)
		}

	}

} // end of main function
