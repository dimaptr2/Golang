package main

import (
	"velkomfood.ru/bengi/bengilib"
)

// Beaglebone's processing of the file and the database loading


// setters and getters

// start point
func main() {

	//gocron.Every(1).Day().At("08:30").Do(runTaskBeagleSrv)

	// test job task every 5 minute
	//	gocron.Every(1).Minute().Do(runTaskBeagleSrv)
	//	fmt.Println("Start beaglebone service")

	// remove, clear and next_run
	//	_, time := gocron.NextRun()
	//	fmt.Println(time)

	//	gocron.Remove(runTaskBeagleSrv)
	//	gocron.Clear()

	// function Start start all the pending jobs
	//	<-gocron.Start()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	//	s := gocron.NewScheduler()
	//	s.Every(3).Seconds().Do(task)
	//	<-s.Start()

	// debug
	bengi.RunTaskBeagleSrv()

}
