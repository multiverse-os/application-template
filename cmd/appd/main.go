package main

import (
	"fmt"
)

func main() {
	fmt.Println("app daemon")
	fmt.Println("==========")
	fmt.Println("review cli tool, it will be similar if its a typical system")
	fmt.Println("service.\n")

	//
	// so load the config values:
	//    [ env => flags => file => default/userinput ]
	//
	//  then use this to load the daemon using whatever server
	//  that is expected.
	//
	//  serviceApp, _ := app.Server.TCP("localhost", 8080)
	//
	//  rcpClient, _ := app.Server.TCP("10.10.10.10", 3000)
	//
	//  then maybe copy data obtained via the rpc and present it over
	// the service connection, or the other way we are providing the
	// service interface. this may be a service for a GUI to use.
	// but the logic for this part of the service should use the library
	// and the daemon logic should be execlusively here.
	//

	fmt.Println("if it is a web application, then this is where we load config")
	fmt.Println("values for loading the web server: the host/address, and the")
	fmt.Println("port.")
	//
	//
	//
	// webApp, _ := app.Server.HTTP("localhost", 8080)
	// webApp.Start()
	//
}
