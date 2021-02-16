package apihandler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// APICall is a type that holds all necessary info for an API call
type APICall struct {
	URL      string
	Method   string
	Header   map[string]string
	Body     []byte
	Insecure bool
}

// Exec executes the underlying API Call and returns the resulting statuscode and error if any occurred
func (a APICall) Exec(i interface{}, client ...http.Client) (statuscode int, err error) {
	if a.Insecure == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// create new Request
	r, err := http.NewRequest(a.Method, a.URL, bytes.NewBuffer(a.Body))
	if err != nil {
		return 0, err
	}
	if a.Header != nil {
		// add header to request r
		for key, value := range a.Header {
			r.Header.Add(key, value)
		}
	}

	resp, err := client[0].Do(r)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return resp.StatusCode, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return resp.StatusCode, err
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}
