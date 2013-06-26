package main

import (
	"flag"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/donovanhide/eventsource"
	"log"
	"scraper"
	"time"
)

var configPath = flag.String("config", "config.json", "The path of the config file")
var delay = flag.Duration("delay", 1*time.Second, "The crawl delay")
var maxVisits = flag.Int("maxvisits", 0, "Maximum number of pages to visit. 0 for unlimited")

func main() {
	flag.Parse()
	srv := eventsource.NewServer()
	configs, err := scraper.Parse(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}
	go scraper.Serve(srv)
	stats := scraper.NewStatistics(srv)
	opts := gocrawl.NewOptions(scraper.NewScraper(stats))
	opts.CrawlDelay = *delay
	opts.LogFlags = gocrawl.LogNone
	opts.MaxVisits = *maxVisits
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err = c.Run(configs); err != gocrawl.ErrMaxVisits {
		log.Fatal(err)
	}
}
