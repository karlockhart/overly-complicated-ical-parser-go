package main

import "github.com/karlockhart/overly-complicated-ical-parser-go/pkg/ical2"
import "log"

import "encoding/json"
import "fmt"

func main() {
	c, err := ical2.ParseIcal2Url("https://calendar.dallasmakerspace.org/events/feed")
	if err != nil {
		log.Fatal(err)
	}

	j, _ := json.Marshal(c)
	fmt.Println(string(j))
}
