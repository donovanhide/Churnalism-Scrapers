package scraper

import (
	"encoding/json"
	"github.com/donovanhide/eventsource"
	"log"
	"net/url"
	"strings"
	"time"
)

type Status int
type LocationMap map[string]map[string]map[string]map[string]int

const (
	Ignored Status = iota
	Visited
	Excluded
	Extracted
	NotFound
	Error
)

var StatusMap = map[Status]string{
	Ignored:   "Ignored",
	Visited:   "Visited",
	Excluded:  "Excluded",
	Extracted: "Extracted",
	NotFound:  "Not Found",
	Error:     "Error",
}

type Stat struct {
	Url    *url.URL
	Status Status
}

type Statistics struct {
	locations LocationMap
}

func (m LocationMap) Add(stat *Stat) {
	sections := strings.Split(stat.Url.Path, "/")
	start, end := strings.Join(sections[0:len(sections)-1], "/")+"/", sections[len(sections)-1]
	host := stat.Url.Host
	status := StatusMap[stat.Status]
	if status == "Extracted" {
		log.Printf("%s%s%s\n", stat.Url.Host, start, end)
	}
	if _, ok := m[host]; !ok {
		m[host] = make(map[string]map[string]map[string]int)
	}
	if _, ok := m[host][status]; !ok {
		m[host][status] = make(map[string]map[string]int)
	}
	if _, ok := m[host][status][start]; !ok {
		m[host][status][start] = make(map[string]int)
	}
	m[host][status][start][end]++
}

func (s *Statistics) Id() string {
	return time.Now().Format(time.RFC3339)
}

func (s *Statistics) Event() string {
	return "summary"
}

func (s *Statistics) Data() string {
	if b, err := json.Marshal(s.locations); err == nil {
		return string(b)
	}
	return ""
}

func NewStatistics(srv *eventsource.Server) chan *Stat {
	c := make(chan *Stat)
	stats := &Statistics{
		locations: make(LocationMap),
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case stat := <-c:
				stats.locations.Add(stat)
			case <-ticker.C:
				srv.Publish([]string{"statistics"}, stats)
			}
		}
	}()
	return c
}
