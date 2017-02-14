package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"strings"
	"time"
)

// Email ...
type Email struct {
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
	sb.WriteString(fmt.Sprintf("body:\n%s\n", e.Body))
	sb.WriteString("\n")

	return sb.String()
}

// Post ...
func (e Emails) Post(u string, delayError int) error {
	log.Printf("Post email to %s", u)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}
	log.Printf("Post json %s", buf.String())

	for {
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
func ParseEmail(bb *bytes.Buffer) (*Email, error) {
	log.Printf("Parse Email: %s", Trunc(bb.String()))
	raw := bb.Bytes()
	m, err := mail.ReadMessage(bb)
	if err != nil {
		return nil, fmt.Errorf("Error parsing email: %s [%s]", Trunc(bb.String()), err)
	}

	header := m.Header

	f := header.Get("From")
	sTo := header.Get("To")
	s := header.Get("Subject")

	d, err := header.Date()
	if err != nil {
		log.Printf("Error parsing date: %s", header.Get("Date"))
	}

	b, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Printf("Error reading email body: %s", err)
	}

	t := strings.Split(sTo, ",")

	return &Email{
		To:      t,
		From:    f,
		Date:    d,
		Subject: s,
		Body:    string(b),
		Raw:     raw,
	}, nil
}
