package lib

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"rsc.io/pdf"
)

func TestDateFormatterSuccess(t *testing.T) {
	tests := []string{
		"Wed, 04 Apr 2012 11:41:06 -0400",
		"Sun, 10 May 2015 09:13:58 -0400",
	}

	for _, test := range tests {
		if d, err := time.Parse(time.RFC1123Z, test); err == nil {
			t.Logf("Success: %v", d)
		} else {
			t.Errorf("Error parsing email date: %s [%s]", test, err)
		}
	}
}

func TestDateFormatterFailure(t *testing.T) {
	tests := []string{
		"Wed, 4 Apr 2012 11:41:06 -0400",
	}

	for _, test := range tests {
		if _, err := time.Parse(time.RFC1123Z, test); err == nil {
			t.Errorf("Error expecting failure: %s", test)
		} else {
			t.Logf("Failure as expected")
		}
	}
}

func TestPdf(t *testing.T) {
	b, err := ioutil.ReadFile("basic.pdf")
	if err != nil {
		t.Fatal(err)
	}
	bb := bytes.NewReader(b)

	r, err := pdf.NewReader(bb, int64(bb.Len()))
	if err != nil {
		t.Fatal(err)
	}

	//log.Print(r.Page(1).V)
	//log.Print(r.Page(1).V.Kind())
	//log.Print(r.Page(1).V.Key("Contents"))

	buf := &bytes.Buffer{}
	reader := r.Page(1).V.Key("Contents").Reader()
	if count, err := buf.ReadFrom(reader); err == nil {
		log.Print(count)
		log.Print(buf.String())

		bufPrint := &bytes.Buffer{}
		for _, v := range buf.Bytes() {
			if strconv.IsPrint(rune(v)) {
				bufPrint.WriteByte(v)
			}
		}
		log.Print(bufPrint)
	} else {
		log.Print(err)
	}
}
