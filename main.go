package apihandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// CallAPI is the raw function that is being used to call the API
func CallAPI(url string, method string, body interface{}, resp interface{}) error {
	// if insecure flag is triggered, don't verify certificates on https
	// if insecure == true {
	// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// }

	log.Println("Calling API with URL:", url)
	client := &http.Client{}

	// setup new request
	var request *http.Request

	// setup a buffer and load body into it, if we are given a body
	if body != nil {
		// convert body to json
		jsonValue, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		request, err = http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			// return fmt.Errorf("couldn't create a request with url %s error was %s", url, err)
			panic(err)
		}
		defer request.Body.Close()

	} else {
		var err error
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			panic(err)
		}
	}

	return getBody(client, request, &resp)
}

func getBody(client *http.Client, request *http.Request, response interface{}) error {
	r, err := client.Do(request)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		panic(err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// return fmt.Errorf("couldn't get a response with url %s error was %s", url, err)
		panic(err)
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		panic(err)
	}
	return nil
}

// GetAPI calls url and writes the answer to the GET request into resp
func GetAPI(url string, resp interface{}) error {
	// log.Println("GetAPI...")
	err := CallAPI(url, http.MethodGet, nil, resp)
	if err != nil {
		panic(err)
	}
	return nil
}

// PostAPI calls url and writes the answer to the POST request with body b into resp
func PostAPI(url string, b interface{}, resp interface{}) error {
	// log.Println("PostAPI...")
	err := CallAPI(url, http.MethodPost, b, &resp)
	if err != nil {
		panic(err)
	}
	return nil
}

// DeleteAPI calls url and write the answer to the DELETE request into resp
func DeleteAPI(url string, body interface{}, resp interface{}) error {
	// log.Println("DeleteAPI...")
	err := CallAPI(url, http.MethodDelete, body, &resp)
	if err != nil {
		panic(err)
	}
	return nil
}
