package lib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// DoGet ...
func DoGet(u string) (*bytes.Buffer, error) {
	log.Printf("DoGet: %s", u)

	url, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url: %s [%s]", u, err)
	}

	client := http.DefaultClient
	if url.Scheme == "https" {
		log.Printf("Using https")
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}}
	}

	resp, err := client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("Error getting url: %s [%s]", u, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting url %s [%d]", u, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", u)
	}

	log.Printf("Response: %d len: %d '%s'", resp.StatusCode, len(body), Trunc(string(body)))
	return bytes.NewBuffer(body), nil
}

// DoPost ...
func DoPost(u string, r io.Reader) error {
	log.Printf("DoPost: %s", u)

	resp, err := http.Post(u, "application/json", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error reading accepted post: %s", err)
		// return err
	}

	log.Printf("Response: %d '%s'", resp.StatusCode, body)
	return nil
}
