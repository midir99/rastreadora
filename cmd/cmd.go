package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/midir99/rastreadora/mpp"
	"github.com/midir99/rastreadora/ws"
	"golang.org/x/net/html"
)

type AlertType string

const (
	AlertTypeGroAlba      = "gro-alba"
	AlertTypeGroAmber     = "gro-amber"
	AlertTypeGroHasVistoA = "gro-hasvistoa"
	AlertTypeMorAmber     = "mor-amber"
	AlertTypeMorCustom    = "mor-custom"
)

func AlertTypesAvailable() []AlertType {
	return []AlertType{
		AlertTypeGroAlba,
		AlertTypeGroAmber,
		AlertTypeGroHasVistoA,
		AlertTypeMorAmber,
		AlertTypeMorCustom,
	}
}

var usageTemplate = `rastreadora is a tool for scraping missing person posters data.

Usage:

    rastreadora [-o output] <alert-type> <from> [until]

Arguments:

    alert-type (string): the type of alerts that you want to collect:{{range .AlertTypes}}
                         - {{.}}{{end}}
    from    (number):    the page number to start scraping missing person posters data.
    until   (number):    the page number to stop scraping missing person posters data, if omitted
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
		AlertTypes []AlertType
	}{AlertTypesAvailable()}
	tmpl := template.Must(template.New("usage").Parse(usageTemplate))
	err := tmpl.Execute(flag.CommandLine.Output(), templateData)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "unable to print help")
	}
}

type Args struct {
	AlertType    AlertType
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
	// Validate the "alert-type" argument
	args.AlertType = AlertType(flag.Arg(0))
	if args.AlertType == "" {
		return nil, fmt.Errorf("<alert-type> argument cannot be empty")
	}
	alertIsValid := false
	for _, s := range AlertTypesAvailable() {
		if args.AlertType == s {
			alertIsValid = true
			break
		}
	}
	if !alertIsValid {
		return nil, fmt.Errorf("\"%s\" is not a valid choice for <alert-type>", args.AlertType)
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
	fmt.Println("rastreadora v0.4.0")
}

func SelectScraperFuncs(alertType AlertType) (func(*html.Node) ([]mpp.MissingPersonPoster, map[int]error), func(uint64) string, error) {
	switch alertType {
	case AlertTypeGroAlba:
		return ws.ScrapeGroAlbaAlerts, ws.MakeGroAlbaUrl, nil
	case AlertTypeGroAmber:
		return ws.ScrapeGroAmberAlerts, ws.MakeGroAmberUrl, nil
	case AlertTypeGroHasVistoA:
		return ws.ScrapeGroHasVistoAAlerts, ws.MakeGroHasVistoAAlertsUrl, nil
	case AlertTypeMorAmber:
		return ws.ScrapeMorAmberAlerts, ws.MakeMorAmberUrl, nil
	case AlertTypeMorCustom:
		return ws.ScrapeMorCustomAlerts, ws.MakeMorCustomUrl, nil
	default:
		return nil, nil, fmt.Errorf("invalid alert-type %v", alertType)
	}
}

func entryLegend(entries int) string {
	if entries == 1 {
		return "entry"
	}
	return "entries"
}

func mppLegend(mpps int) string {
	if mpps == 1 {
		return "missing person poster"
	}
	return "missing person posters"
}

func Scrape(pageUrl string, scraper func(*html.Node) ([]mpp.MissingPersonPoster, map[int]error), ch chan []mpp.MissingPersonPoster) {
	doc, err := ws.RetrieveDocument(pageUrl)
	if err != nil {
		// log.Printf("Done processing %s; 0 entries collected; %s", pageUrl, err)
		log.Printf("0 entries collected from %s; %s", pageUrl, err)
		ch <- []mpp.MissingPersonPoster{}
		return
	}
	mpps, errs := scraper(doc)
	mppsLen := len(mpps)
	entryWord := entryLegend(mppsLen)
	if errsLen := len(errs); errsLen > 0 {
		messages := []string{}
		for entryNumber, err := range errs {
			messages = append(messages, fmt.Sprintf("entry #%d: %s", entryNumber, err))
		}
		message := strings.Join(messages, ",")
		// log.Printf("Done processing %s; %d %s collected; unable to retrieve %d, details: %s", pageUrl, mppsLen, entryWord, errsLen, message)
		log.Printf("%d %s collected from %s; unable to retrieve %d, details: %s", mppsLen, entryWord, pageUrl, errsLen, message)
	} else {
		// log.Printf("Done processing %s; %d %s collected", pageUrl, mppsLen, entryWord)
		log.Printf("%d %s collected from %s", mppsLen, entryWord, pageUrl)
	}
	ch <- mpps
}

func Execute(args *Args) {
	if args.PrintVersion {
		PrintVersion()
		os.Exit(0)
	}
	scraper, makeUrl, err := SelectScraperFuncs(args.AlertType)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	ch := make(chan []mpp.MissingPersonPoster)
	for pageNum := args.PageFrom; pageNum <= args.PageUntil; pageNum++ {
		pageUrl := makeUrl(pageNum)
		go Scrape(pageUrl, scraper, ch)
	}
	mpps := []mpp.MissingPersonPoster{}
	pagesCount := args.PageUntil - args.PageFrom + 1
	for curPage := uint64(1); curPage <= pagesCount; curPage++ {
		mpps = append(mpps, <-ch...)
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
	mppsLen := len(mpps)
	mppWord := mppLegend(mppsLen)
	log.Printf("%d %s collected", mppsLen, mppWord)
}
