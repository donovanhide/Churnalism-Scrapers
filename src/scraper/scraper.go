package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"log"
	// "github.com/bcampbell/arts"
	// "io/ioutil"
	"net/http"
	"regexp"
)

var exclusions = regexp.MustCompile(`\.(?i:jpg|gif|doc|docx|pdf|zip|ppt|xls|png)$`)

type Scraper struct {
	gocrawl.DefaultExtender
	configs map[string]*Config
	stats   chan *Stat
}

func NewScraper(c chan *Stat) *Scraper {
	return &Scraper{
		stats: c,
	}
}

func (s *Scraper) Start(configs interface{}) interface{} {
	s.configs = make(map[string]*Config)
	for _, v := range configs.(gocrawl.S) {
		config := v.(*Config)
		s.configs[config.Host] = config
	}
	return configs
}

func (s *Scraper) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	defer func() {
		if r := recover(); r != nil {
			s.Log(gocrawl.LogError, gocrawl.LogError, fmt.Sprint("Visit Recovered in ", r))
		}
	}()
	url := ctx.NormalizedURL()
	// log.Println("Visit:", url.String(), s.configs[url.Host].Extractable(url.String()))
	if !s.configs[url.Host].Extractable(url.String()) {
		s.stats <- &Stat{url, Visited}
		return nil, true
	}
	s.stats <- &Stat{url, Extracted}
	// html, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, true
	// }
	// article, err := arts.Extract(html, ctx.URL().String(), false)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, true
	// }
	// log.Println(ctx.NormalizedURL().String(), article.Headline)
	return nil, true
}

func (s *Scraper) Error(err *gocrawl.CrawlError) {
	log.Println(err, err.Ctx.URL())
	switch err.Error() {
	case "404 Not found":
		s.stats <- &Stat{err.Ctx.URL(), NotFound}
	case "500 Error":
		s.stats <- &Stat{err.Ctx.URL(), Error}
	default:
		s.stats <- &Stat{err.Ctx.URL(), Error}
	}
}

func (s *Scraper) Filter(ctx *gocrawl.URLContext, isVisited bool) (visit bool) {
	defer func() {
		if r := recover(); r != nil {
			s.Log(gocrawl.LogError, gocrawl.LogError, fmt.Sprint("Filter Recovered in ", r))
		}
	}()
	url := ctx.NormalizedURL()
	config, ok := s.configs[url.Host]
	visit = !isVisited && !exclusions.MatchString(url.String()) && ok && config.Visitable(url.String())
	if !visit && ok {
		s.stats <- &Stat{url, Ignored}
	}
	return
}
