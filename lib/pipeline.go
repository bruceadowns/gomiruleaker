package lib

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"
)

// Generate ...
// * enumerates urls from Targets
// * retrieves the http.Get
// * send bytes it along to the parse channel
func Generate(t *Targets, d int) chan *bytes.Buffer {
	log.Printf("Generate sources: %s [%d]", t, d)

	out := make(chan *bytes.Buffer)

	go func() {
		defer close(out)

		for idx := t.Start; idx <= t.End; idx++ {
			if bb, err := DoGet(fmt.Sprintf("%s%d", t.Prefix, idx)); err == nil {
				out <- bb
			} else {
				break
			}

			if d > 0 {
				log.Printf("Sleep %dms between http gets", d)
				time.Sleep(time.Duration(d) * time.Millisecond)
			}
		}
	}()

	return out
}

// Parse ...
// * takes bytes channel input
// * normalizes input to text
// * parses email
// * sends email to accumulator
func Parse(in chan *bytes.Buffer, c int) chan *Email {
	log.Printf("Parser Count: %d", c)

	out := make(chan *Email)
	var wg sync.WaitGroup

	for i := 0; i < c; i++ {
		log.Printf("Start Parser %d", i)
		wg.Add(1)

		go func() {
			for bb := range in {
				if m, err := ParseEmail(bb); err == nil {
					out <- m
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
func Accum(in chan *Email, b int) chan Emails {
	log.Printf("Accum channel in batches of %d", b)

	out := make(chan Emails)

	go func() {
		defer close(out)

		emails := make(Emails, 0)
		for email := range in {
			log.Printf("Accumulate email: %s", email)

			emails = append(emails, email)
			if len(emails) >= b {
				log.Printf("Send %d emails to posters", len(emails))
				out <- emails
				emails = make(Emails, 0)
			}
		}

		if len(emails) > 0 {
			log.Printf("Send %d emails to posters", len(emails))
			out <- emails
		}
	}()

	return out
}

// Post ...
func Post(in chan Emails, u string, d int) chan int {
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
				if err := emails.Post(u, d); err != nil {
					log.Printf("Error posting: %s", err)
				}
			}

			out <- len(emails)
		}
	}()

	return out
}
