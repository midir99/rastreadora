package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"

	"github.com/midir99/rastreadora/mpp"
	"github.com/midir99/rastreadora/ws"
	"golang.org/x/net/html"
)

type Scraper string

const (
	ScraperGroAlba      = "gro-alba"
	ScraperGroAmber     = "gro-amber"
	ScraperGroHasVistoA = "gro-hasvistoa"
	ScraperMorAmber     = "mor-amber"
	ScraperMorCustom    = "mor-custom"
)

func ScrapersAvailable() []Scraper {
	return []Scraper{
		ScraperGroAlba,
		ScraperGroAmber,
		ScraperGroHasVistoA,
		ScraperMorAmber,
		ScraperMorCustom,
	}
}

var usageTemplate = `rastreadora is a tool for scraping missing person posters data.

Usage:

    rastreadora [-o output] <scraper> <from> [until]

Arguments:

    scraper (string): the scraper that will be used to extract data, available values:{{range .Scrapers}}
                      - {{.}}{{end}}
    from    (number): the page number to start scraping missing person posters data.
    until   (number): the page number to stop scraping missing person posters data, if omitted
                      the program will only scrap data from the page number specified by the
                      <from> argument.

Flags:

    -o      (string): the filename where the data will be stored, if omitted the data will be
                      dumped in STDOUT.
    -V      (bool):   print the version of the program.
    -h      (bool):   print this usage message.
`

func Usage() {
	templateData := struct {
		Scrapers []Scraper
	}{Scrapers: ScrapersAvailable()}
	tmpl := template.Must(template.New("usage").Parse(usageTemplate))
	err := tmpl.Execute(flag.CommandLine.Output(), templateData)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "unable to print help")
	}
}

type Args struct {
	Scraper      Scraper
	PageFrom     uint64
	PageUntil    uint64
	Out          string
	SkipCert     bool
	PrintVersion bool
}

func ParseArgs() (*Args, error) {
	args := Args{}
	flag.StringVar(&args.Out, "o", "", "the filename where the data will be stored, if omitted the data will be dumped in STDOUT.")
	flag.BoolVar(&args.SkipCert, "scert", false, "skip the verification of the server's certificate chain and hostname.")
	flag.BoolVar(&args.PrintVersion, "V", false, "print the version of the program.")
	flag.Usage = Usage
	flag.Parse()
	if args.PrintVersion {
		return &args, nil
	}
	// Validate the "scraper" argument
	args.Scraper = Scraper(flag.Arg(0))
	if args.Scraper == "" {
		return nil, fmt.Errorf("<scraper> argument cannot be empty")
	}
	scraperIsValid := false
	for _, s := range ScrapersAvailable() {
		if args.Scraper == s {
			scraperIsValid = true
			break
		}
	}
	if !scraperIsValid {
		return nil, fmt.Errorf("\"%s\" is not a valid choice for <scraper>", args.Scraper)
	}
	// Validate the "from" argument
	if flag.Arg(1) == "" {
		return nil, fmt.Errorf("<from> argument cannot be empty")
	}
	pF, err := strconv.ParseUint(flag.Arg(1), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("\"%s\" is not a valid number for <from>", flag.Arg(1))
	}
	args.PageFrom = pF
	// Validate the "until" argument
	if flag.Arg(2) == "" {
		args.PageUntil = args.PageFrom
	} else {
		pU, err := strconv.ParseUint(flag.Arg(2), 10, 0)
		if err != nil {
			return nil, fmt.Errorf("\"%s\" is not a valid number for [until]", flag.Arg(2))
		}
		args.PageUntil = pU
	}
	// Validate "from" value is lower or equal to "until" value
	if args.PageFrom > args.PageUntil {
		return nil, fmt.Errorf("<from> value must be lower or equal to [until] value")
	}
	return &args, nil
}

func PrintVersion() {
	fmt.Println("rastreadora v0.3.0	")
}

func SelectScraperFuncs(scraper Scraper) (func(*html.Node) []mpp.MissingPersonPoster, func(uint64) string, error) {
	switch scraper {
	case ScraperGroAlba:
		return ws.ScrapeGroAlbaAlerts, ws.MakeGroAlbaUrl, nil
	case ScraperGroAmber:
		return ws.ScrapeGroAmberAlerts, ws.MakeGroAmberUrl, nil
	case ScraperGroHasVistoA:
		return ws.ScrapeGroHasVistoAAlerts, ws.MakeGroHasVistoAAlertsUrl, nil
	case ScraperMorAmber:
		return ws.ScrapeMorAmberAlerts, ws.MakeMorAmberUrl, nil
	case ScraperMorCustom:
		return ws.ScrapeMorCustomAlerts, ws.MakeMorCustomUrl, nil
	default:
		return nil, nil, fmt.Errorf("invalid scraper %v", scraper)
	}
}

func Execute(args *Args) {
	if args.PrintVersion {
		PrintVersion()
		os.Exit(0)
	}
	scraper, makeUrl, err := SelectScraperFuncs(args.Scraper)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	ch := make(chan []mpp.MissingPersonPoster)
	for pageNum := args.PageFrom; pageNum <= args.PageUntil; pageNum++ {
		pageUrl := makeUrl(pageNum)
		log.Printf("Processing %s ...\n", pageUrl)
		go ws.Scrape(pageUrl, scraper, ch)
	}
	mpps := []mpp.MissingPersonPoster{}
	pagesCount := args.PageUntil - args.PageFrom + 1
	for curPage := uint64(1); curPage <= pagesCount; curPage++ {
		mpps = append(mpps, <-ch...)
		log.Printf("%d out of %d page(s) have been scraped\n", curPage, pagesCount)
	}
	output, err := json.Marshal(mpps)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	if args.Out != "" {
		if os.WriteFile(args.Out, output, 0664) != nil {
			log.Fatalf("Error: %s", err)
		}
	} else {
		_, err := os.Stdout.Write(output)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
	log.Printf("%d missing person poster(s) were processed\n", len(mpps))
}
