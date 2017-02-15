package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang/snappy"
)

// Email ...
type Email struct {
	ID      int
	Type    string
	To      []string
	From    string
	Date    time.Time
	Subject string
	Body    string
	Raw     []byte
}

// Emails ...
type Emails []*Email

func (e *Email) String() string {
	sb := &bytes.Buffer{}

	sb.WriteString("\n")

	sb.WriteString("to: ")
	for k, v := range e.To {
		if k > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s", v))
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("from: %s\n", e.From))
	sb.WriteString(fmt.Sprintf("date: %s\n", e.Date.Format(time.RFC1123Z)))
	sb.WriteString(fmt.Sprintf("subject: %s\n", e.Subject))
	sb.WriteString(fmt.Sprintf("body:\n%s\n", Trunc(e.Body)))
	sb.WriteString("\n")

	return sb.String()
}

// Post ...
func (e Emails) Post(u string, delayError int, compress bool) error {
	log.Printf("Post email to %s", u)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}
	log.Printf("Post %d json bytes: %s", buf.Len(), Trunc(buf.String()))

	for {
		if compress {
			sb := snappy.Encode(nil, buf.Bytes())
			buf = bytes.NewBuffer(sb)
			log.Printf("Post %d compressed json bytes", buf.Len())
		}

		r := bytes.NewReader(buf.Bytes())
		if err := DoPost(u, r); err == nil {
			break
		} else {
			log.Printf("Error posting to miru-leaks: %s", err)
		}

		if delayError > 0 {
			log.Printf("Miru delay on error %dms", delayError)
			time.Sleep(time.Duration(delayError) * time.Millisecond)
		}
	}

	return nil
}

// ParseEmail ...
func ParseEmail(e *Email) error {
	log.Printf("Parse Email: %s", Trunc(string(e.Raw)))

	m, err := mail.ReadMessage(bytes.NewBuffer(e.Raw))
	if err != nil {
		return fmt.Errorf("Error parsing email: %s [%s]", Trunc(string(e.Raw)), err)
	}

	header := m.Header

	e.From = header.Get("From")
	e.To = strings.Split(header.Get("To"), ",")
	e.Subject = header.Get("Subject")

	d, err := header.Date()
	if err != nil {
		log.Printf("Error parsing date: %s", header.Get("Date"))
	}
	e.Date = d

	b, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Printf("Error reading email body: %s", err)
	}

	e.Body = string(b)

	return nil
}

// Save ...
func (e Email) Save(d string) error {
	f := fmt.Sprintf("%s_%d.eml", e.Type, e.ID)
	if d == "" {
		log.Printf("Not saving email: %s", f)
	} else {
		ff := fmt.Sprintf("%s/%s", d, f)
		log.Printf("Save Email to %s [%s]", ff, Trunc(e.String()))

		b, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("Error encoding email: %s", err)
		}

		ioutil.WriteFile(ff, b, os.ModePerm)
	}

	return nil
}
