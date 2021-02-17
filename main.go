package apihandler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
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
func (a APICall) Exec(i interface{}) (statuscode int, err error) {
	client := httpClient()

	if a.Insecure == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// create new Request
	r, err := http.NewRequest(a.Method, a.URL, bytes.NewBuffer(a.Body))
	if err != nil {
		return 0, err
	}

	if a.Cookie != nil {
		r.AddCookie(a.Cookie)
	}

	if a.Header != nil {
		// add header to request r
		for key, value := range a.Header {
			r.Header.Add(key, value)
		}
	}

	resp, err := client.Do(r)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return resp.StatusCode, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(data, &i)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}
