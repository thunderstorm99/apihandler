package apihandler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APICall is a type that holds all necessary info for an API call
type APICall struct {
	URL      string
	Method   string
	Header   map[string]string
	Body     interface{}
	Insecure bool
}

// Exec executes the underlying API Call
func (a APICall) Exec(i interface{}) error {
	if a.Insecure == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	r, err := http.NewRequest(a.Method, a.URL, nil)
	if err != nil {
		return err
	}
	if a.Header != nil {
		// add header to request r
		for key, value := range a.Header {
			r.Header.Add(key, value)
		}
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return err
	}

	if resp.StatusCode != 200 {
		// non standard Status code returns
		return fmt.Errorf("HTTP statuscode is not 200, error is: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return err
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	return nil
}
