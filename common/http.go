package common

import (
	"bytes"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Headers is a simple structure to list key/value pair on http request header.
type Headers map[string]string

// Get data from a API service and decode json value automatically.
func Get(url string, data interface{}) (response *http.Response, err error) {
	response, err = http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(data)

	return
}

// Post to an API endpoint, json encoding and decoding req and data respectivelly.
func Post(url string, req interface{}, data interface{}) (response *http.Response, err error) {
	var reqBody []byte

	reqBody, err = json.Marshal(req)
	response, err = http.Post(url, "test/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(data)

	return
}

// Request execute a request with headers and method setup
func Request(method string, headers Headers, url string, req interface{}, data interface{}) (response *http.Response, err error) {
	var request *http.Request

	if method == "POST" || method == "PUT" {
		var reqBody []byte
		// Encode the request
		reqBody, err = json.Marshal(req)
		if err != nil {
			return
		}
		// Create the request
		request, err = http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}

	// Set Request header from headers
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	//fmt.Printf("request:\n%s\n", request)
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return response, fmt.Errorf("%s", response.Status)
	}

	defer response.Body.Close()

	//err = json.NewDecoder(response.Body).Decode(data)

	buf, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(buf, data)

	return
}

// SetQuery creates a Query URL given with path args
// Returns the query part (to set in Request.URL.RawQuery)
func SetQuery(urlInput *url.URL, args map[string]string) (ret string) {
	if urlInput == nil {
		return
	}

	if len(args) > 0 {
		q := url.Values{}
		for key, value := range args {
			q.Add(key, value)

		}
		urlInput.RawQuery = q.Encode()
	}
	return
}

// UpdateQuery update an existing Query URL given with path args
// Returns the query part (to set in Request.URL.RawQuery)
func UpdateQuery(urlInput *url.URL, args map[string]string) (ret string) {
	if urlInput == nil {
		return
	}

	if len(args) > 0 {
		q := urlInput.Query()
		for key, value := range args {
			q.Set(key, value)

		}
		urlInput.RawQuery = q.Encode()
	}
	return
}
