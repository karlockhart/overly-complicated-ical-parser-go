package main

import "github.com/karlockhart/overly-complicated-ical-parser-go/pkg/ical2"

func main() {
	ical2.ParseIcal2Url("https://calendar.dallasmakerspace.org/events/feed")
}
