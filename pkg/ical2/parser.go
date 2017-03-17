package ical2

import "net/http"

import "io/ioutil"

import "strings"
import "time"

const dateLayout string = "20060102T150405Z"

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
	StartDate   time.Time
	EndDate     time.Time
	DateStamp   time.Time
	UID         string
	Summary     string
	Description string
	Location    string
	URL         string
}

func parseTime(in string) time.Time {
	t, err := time.Parse(dateLayout, in)
	if err != nil {
		panic(err)
	}

	return t
}

func (e *Event) populate(t string, val string) {
	switch t {
	case "DTSTART":
		e.StartDate = parseTime(val)
	case "DTEND":
		e.EndDate = parseTime(val)
	case "DTSTAMP":
		e.DateStamp = parseTime(val)
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
	c.children = make(map[string]*node)
	p.children[pat] = c
	return c
}

func (p *node) addExistingChild(pat string, n *node) *node {
	p.children[pat] = n
	return n
}

func splitElement(in string) string {
	parts := strings.Split(in, ":")
	if len(parts) > 1 {
		return parts[1]
	}

	panic("Could not parse value.")
}

func createEvent(c *Calendar, t string, val string) *Calendar {
	e := new(Event)
	c.Events = append(c.Events, e)
	c.curr = e
	return c
}

func noOp(c *Calendar, t string, val string) *Calendar {
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

func initialize() *node {
	root := new(node)
	root.children = make(map[string]*node)
	curr := root.addChild("BEGIN:VCALENDAR", nil, nil)
	start := *curr
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
	curr = curr.addChild("END:VEVENT", nil, noOp)
	curr.addExistingChild("BEGIN:VEVENT", curr.prev)
	curr = curr.addChild("END:VCALENDAR", nil, noOp)

	return &start
}

func preprocessICal2String(in string) string {
	return strings.TrimSpace(strings.Replace(in, "\\n", " ", -1))
}

func trimToFirstDirective(in string) string {
	idx := strings.Index(in, "BEGIN:VCALENDAR")
	return in[idx : len(in)-1]
}

// ParseIcal2Url parses an ICal2 url into a Calendar.
func ParseIcal2Url(url string) (*Calendar, error) {
	curr := initialize()

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	s := trimToFirstDirective(preprocessICal2String(string(bytes)))
	lines := strings.Split(s, "\n")
	c := new(Calendar)

	for _, l := range lines {
		s := strings.TrimSpace(l)
		// This is a state transition.
		if _, ok := curr.children[s]; ok {
			curr.children[s].populate(c, s, s)
			curr = curr.children[s]
			continue
		}

		parts := strings.Split(s, ":")
		// This is NOT a state transition.
		if _, ok := curr.children[parts[0]]; ok {
			// Run extractors.
			for _, e := range curr.children[parts[0]].extractor {
				s = e(s)
			}

			// Populate the Calendar.
			curr.children[parts[0]].populate(c, parts[0], s)
			continue
		}

	}

	return c, err
}
