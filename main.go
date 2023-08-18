package main

import (
	"smallDFS/dataserver"
	"smallDFS/nameserver"
)

func runNameserver() {
	ns := nameserver.New(3, 3)
	ns.Add("http://localhost:8000")
	ns.Run()
}

func runDataServer() {
	ds := dataserver.New("8000", "ds-1")
	ds.Run()
}

func main() {
	runDataServer()
	//runNameserver()
}
