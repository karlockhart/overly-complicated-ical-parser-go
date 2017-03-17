package ical2

import "net/http"

import "io/ioutil"
import "log"

type finalizer interface {
	finalize(in string)
}

type node struct {
	children  map[string]*node
	extractor []func(in string) string
	finalize  func(f finalizer, in string)
	prev      *node
}

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
	curr = curr.addChild("BEGIN:VEVENT", nill)

}

type populatable interface {
	populate(string, string)
}

// Parse stuff.
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
