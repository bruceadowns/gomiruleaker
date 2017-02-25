package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	wikiLeaksURLTemplate    = "https://www.wikileaks.org/%s/get/%d"
	foiaMetaDataURLTemplate = "https://foia.state.gov/searchapp/Search/SubmitSimpleQuery?searchText=*&beginDate=false&endDate=false&collectionMatch=%s&postedBeginDate=false&postedEndDate=false&caseNumber=false&page=1&start=%d&limit=%d"
	foiaDocumentURLTemplate = "https://foia.state.gov/searchapp/%s"
)

var (
	reDate = regexp.MustCompile(`new Date\(([0-9]+)\)`)
)

func generateWikiLeaks(t *InputTarget, out chan *Email, delay int, wg *sync.WaitGroup) error {
	log.Printf("Generate wikileak sources")
	defer wg.Done()

	for idx := t.Start; idx <= t.End; idx++ {
		u := fmt.Sprintf(wikiLeaksURLTemplate, t.SubType, idx)
		if bb, err := DoGet(u); err == nil {
			out <- &Email{
				Type: t.SubType,
				ID:   strconv.Itoa(idx),
				Raw:  bb.Bytes(),
			}
		} else {
			log.Printf("Error occurred getting %s", u)
			// should return error when not 404
			break
		}

		if delay > 0 {
			log.Printf("Sleep %dms between http gets", delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	log.Printf("Done generating wikileak sources")
	return nil
}

func generateFoia(t *InputTarget, out chan *Email, delay int, wg *sync.WaitGroup) error {
	log.Printf("Generate foia sources")
	defer wg.Done()

	start := t.Start
	end := t.End
	limit := t.Limit

	if start+limit > end {
		limit = end - start + 1
	}

	var urls []string
	for {
		u := fmt.Sprintf(foiaMetaDataURLTemplate, t.SubType, start, limit)
		if bb, err := DoGet(u); err == nil {
			bbClean := reDate.ReplaceAll(bb.Bytes(), []byte("$1"))

			var results FoiaResults
			if derr := json.Unmarshal(bbClean, &results); derr != nil {
				log.Printf("Error decoding foia metadata: %s", derr)
				break
			}

			for _, v := range results.Results {
				urls = append(urls, fmt.Sprintf(foiaDocumentURLTemplate, v.PdfLink))
			}

			start = start + limit
			if end > results.TotalHits {
				end = results.TotalHits
			}
			if limit > end-start+1 {
				limit = end - start + 1
			}

			if limit < 1 {
				break
			}
		} else {
			log.Printf("Error getting foia metadata: %s", err)
			break
		}
	}

	for _, v := range urls {
		if bb, err := DoGet(v); err == nil {
			out <- &Email{
				Type: t.SubType,
				ID:   path.Base(v),
				Raw:  bb.Bytes(),
			}
		} else {
			log.Printf("Error occurred getting %s", v)
			break
		}

		if delay > 0 {
			log.Printf("Sleep %dms between http gets", delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	log.Printf("Done generating foia sources")
	return nil
}

// Generate ...
// * enumerates urls from Targets
// * retrieves the http.Get
// * send bytes it along to the parse channel
func Generate(in []*InputTarget, delay int) chan *Email {
	log.Printf("Generate sources")

	out := make(chan *Email)
	var wg sync.WaitGroup

	for _, v := range in {
		switch v.Type {
		case "wikiLeaks":
			wg.Add(1)
			go generateWikiLeaks(v, out, delay, &wg)
		case "foia":
			wg.Add(1)
			go generateFoia(v, out, delay, &wg)
		default:
			log.Printf("Unknown target type")
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Parse ...
// * takes bytes channel input
// * normalizes input to text
// * parses email
// * sends email to accumulator
func Parse(in chan *Email, c int) chan *Email {
	log.Printf("Parser Count: %d", c)

	out := make(chan *Email)
	var wg sync.WaitGroup

	for i := 0; i < c; i++ {
		log.Printf("Start Parser %d", i+1)
		wg.Add(1)

		go func() {
			for rawemail := range in {
				if err := rawemail.Parse(); err == nil {
					out <- rawemail
				} else {
					log.Print(err)
				}
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Accum ...
func Accum(in chan *Email, b int, d string) (chan Emails, error) {
	log.Printf("Accum channel in batches of %d", b)

	out := make(chan Emails)

	if len(d) > 0 {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return nil, err
		}
	}

	go func() {
		defer close(out)

		emails := make(Emails, 0)
		for email := range in {
			log.Printf("Accumulate email: %s", email)

			email.Save(d)
			emails = append(emails, email)

			if len(emails) >= b {
				log.Printf("Send %d emails to poster", len(emails))
				out <- emails
				emails = make(Emails, 0)
			}
		}

		if len(emails) > 0 {
			log.Printf("Send %d emails to poster", len(emails))
			out <- emails
		}
	}()

	return out, nil
}

// Post ...
func Post(in chan Emails, u string, d int, c bool) chan int {
	if u == "" {
		u = "stdout"
	}
	log.Printf("Post channel to %s [%dms]", u, d)

	out := make(chan int)

	go func() {
		defer close(out)

		for emails := range in {
			if u == "stdout" {
				for k, v := range emails {
					log.Printf("Post to stdout: %d: %s", k, Trunc(v.String()))
				}
			} else {
				if err := emails.Post(u, d, c); err != nil {
					log.Printf("Error posting: %s", err)
				}
			}

			out <- len(emails)
		}
	}()

	return out
}
