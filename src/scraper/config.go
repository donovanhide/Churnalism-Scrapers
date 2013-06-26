package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/purell"
	"net/url"
	"os"
	"regexp"
)

type Config struct {
	Name, Category, Seed, Visit, Extract, Host string
	DocType                                    uint32
	visitable, extractable                     *regexp.Regexp
}

func Parse(path string) (config gocrawl.S, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	var s []Config
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&s); err != nil {
		return
	}
	config = make(gocrawl.S)
	for i := range s {
		if s[i].visitable, err = regexp.Compile(s[i].Visit); err != nil {
			return
		}
		if s[i].extractable, err = regexp.Compile(s[i].Extract); err != nil {
			return
		}
		seed := purell.MustNormalizeURLString(s[i].Seed, purell.FlagsAllGreedy)
		var u *url.URL
		if u, err = url.Parse(seed); err != nil {
			return
		}
		s[i].Host = u.Host
		config[s[i].Seed] = &s[i]
	}
	return
}

func (c *Config) Visitable(url string) bool {
	return c.visitable.MatchString(url)
}

func (c *Config) Extractable(url string) bool {
	return c.extractable.MatchString(url)
}
