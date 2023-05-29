package apihandler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
)

// APICall is a type that holds all necessary info for an API call
type APICall struct {
	URL      string
	Method   string
	Header   map[string]string
	Body     []byte
	Insecure bool
	Cookie   *http.Cookie
}

func httpClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Jar: jar}
	return client
}

// Exec executes the underlying API Call and returns the resulting statuscode and error if any occurred
func (a APICall) Exec(i any) (statuscode int, errormessage any, err error) {
	client := httpClient()

	if a.Insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// create new Request
	r, err := http.NewRequest(a.Method, a.URL, bytes.NewBuffer(a.Body))
	if err != nil {
		return 0, nil, err
	}

	// add cookie if needed
	if a.Cookie != nil {
		r.AddCookie(a.Cookie)
	}

	if a.Header != nil {
		// add header to request r
		for key, value := range a.Header {
			r.Header.Add(key, value)
		}
	}

	// execute request
	resp, err := client.Do(r)
	if err != nil {
		return 0, nil, err
	}

	// read body to data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	defer resp.Body.Close()

	// unmarshal data onto i (the given data structure)
	err = json.Unmarshal(data, &i)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	// also unmarshal onto errormessage (in case this is the error message from the server)
	err = json.Unmarshal(data, &errormessage)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, errormessage, nil
}
