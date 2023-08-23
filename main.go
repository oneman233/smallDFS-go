package main

import (
	"smallDFS/dataserver"
	"smallDFS/nameserver"
)

func runNameserver() {
	ns := nameserver.New(3, 3, 5)
	ns.Add("http://localhost:8000")
	ns.Run()
}

func runDataServer() {
	ds := dataserver.New("8000", "ds-1")
	ds.Run()
}

// put a.txt b/a.txt
// read b/a.txt
// dump
// undump
// tree
func main() {
	//runDataServer()
	runNameserver()
}
