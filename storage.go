package main

import "strconv"

func makeIdGen() func() string {
	var nextID = 0
	return func() string {
		nextID++
		return strconv.Itoa(nextID)
	}
}

var nextID = makeIdGen()
var db = Todos{}

