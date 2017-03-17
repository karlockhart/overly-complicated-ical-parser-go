package ical2

import "net/http"

import "io/ioutil"
import "log"

type node struct {
	children  map[string]*node
	extractor []func(in string) string
	populate  func(c *Calendar, t string, val string) *Calendar
	prev      *node
}

// Calendar represents an ICal2 Calendar.
type Calendar struct {
	Version  string
	ProdID   string
	CalScale string
	Events   []*Event
	curr     *Event
}

func (c *Calendar) populate(t string, val string) {
	switch t {
	case "VERSION":
		c.Version = val
	case "PRODID":
		c.ProdID = val
	case "CALSCALE":
		c.CalScale = val
	}
}

// Event represents an ICal2 Event.
type Event struct {
	StartDate   int64
	EndDate     int64
	DateStamp   int64
	UID         string
	Summary     string
	Description string
	Location    string
	URL         string
}

func (e *Event) populate(t string, val string) {
	switch t {
	case "DTSTART":
		e.StartDate = 111
	case "DTEND":
		e.EndDate = 111
	case "DTSTAMP":
		e.DateStamp = 111
	case "UID":
		e.UID = val
	case "SUMMARY":
		e.Summary = val
	case "DESCRIPTION":
		e.Description = val
	case "LOCATION":
		e.Location = val
	case "URL":
		e.UID = val
	}
}

type calendar struct {
	calInfo map[string]string
	events  []map[string]string
}

type parserObject struct {
	kv map[string]string
}

func (p *node) addChild(pat string, ext []func(in string) string, pFunc func(c *Calendar, t string, val string) *Calendar) *node {
	c := new(node)
	c.prev = p
	c.extractor = ext
	c.populate = pFunc
	p.children[pat] = c
	return c
}

func splitElement(in string) string {
	return "blah"
}

func createEvent(c *Calendar, t string, val string) *Calendar {
	e := new(Event)
	c.Events = append(c.Events, e)
	c.curr = e
	return c
}

func populateCalendar(c *Calendar, t string, val string) *Calendar {
	c.populate(t, val)
	return c
}

func populateEvent(c *Calendar, t string, val string) *Calendar {
	c.curr.populate(t, val)
	return c
}

func initialize() {
	root := new(node)
	root.children = make(map[string]*node)
	curr := root.addChild("BEGIN:VCALENDAR", nil, nil)
	splitOnlyChain := []func(in string) string{splitElement}
	curr.addChild("VERSION", splitOnlyChain, populateCalendar)
	curr.addChild("PRODID", splitOnlyChain, populateCalendar)
	curr.addChild("CALSCALE", splitOnlyChain, populateCalendar)
	curr = curr.addChild("BEGIN:VEVENT", nil, createEvent)
	curr.addChild("DTSTART", splitOnlyChain, populateEvent)
	curr.addChild("DTSTAMP", splitOnlyChain, populateEvent)
	curr.addChild("UID", splitOnlyChain, populateEvent)
	curr.addChild("SUMMARY", splitOnlyChain, populateEvent)
	curr.addChild("DESCRIPTION", splitOnlyChain, populateEvent)
	curr.addChild("LOCATION", splitOnlyChain, populateEvent)
	curr.addChild("URL", splitOnlyChain, populateEvent)
	curr = curr.addChild("END:VEVENT", nil, createEvent)
	curr = curr.addChild("END:VCALENDAR", nil, nil)
}

type populatable interface {
	populate(string, string)
}

// ParseIcal2Url parses an ICal2 url into a Calendar.
func ParseIcal2Url(url string) error {
	initialize()

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	log.Println(string(bytes))

	return nil
}
