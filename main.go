package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

const host = "https://github.com"

var (
	topN  int
	lang  string
	since string
)

type RepsList struct {
	URL         string
	Description string
}

func main() {
	fmt.Print(
		`-----------------------------------------------`,
		"\n",
		"CLI tool Github trending scraper",
		"\n",
	)

	flag.IntVar(&topN, "top_n", 5, "Top n repositories")
	flag.StringVar(&lang, "lang", "go", "Programming language")
	flag.StringVar(&since, "since", "weekly", "Date range")
	flag.Parse()

	if topN > 25 {
		fmt.Print(
			"\n",
			"---Maximum 25 repositories allowed---",
			"\n",
			"---25 repositories will be shown---",
			"\n", "\n",
		)
		topN = 25
	}

	invalidSince := since != "daily" && since != "weekly" && since != "monthly"

	if invalidSince {
		fmt.Print(
			"\n",
			"---Invalid date range---",
			"\n",
			"---Date range will be weekly---",
			"\n", "\n",
		)
		since = "weekly"
	}

	c := colly.NewCollector()

	list := make([]RepsList, 0, topN)

	c.OnHTML("h1", func(h *colly.HTMLElement) {
		if h.Attr("class") == "h3 lh-condensed" {
			if len(list) <= topN-1 {
				list = append(list, RepsList{
					URL:         host + h.ChildAttr("a", "href"),
					Description: h.DOM.Next().Text(),
				})
			}
		}
	})

	url := fmt.Sprintf("%v/trending/%v?since=%v", host, lang, since)

	c.Visit(url)

	fmt.Print(
		`-----------------------------------------------`,
		"\n",
		"Top "+strconv.Itoa(topN)+" repositories",
		"\n",
		"Language: "+lang,
		"\n",
		"Date range: "+since,
		"\n",
		`-----------------------------------------------`,
		"\n",
	)

	for ind, rep := range list {
		fmt.Print(
			strconv.Itoa(ind+1)+". ",
			rep.URL,
			"\n",
			strings.TrimSpace(rep.Description),
			"\n", "\n",
		)
	}
}
