package wso2requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

// MakeWSO2Request makes a generic WSO2 request
// toReturn should be a pointer
func MakeWSO2Request(method, url string, body interface{}, toReturn interface{}) *nerr.E {
	return MakeWSO2RequestWithHeaders(method, url, body, toReturn, nil)
}

// MakeWSO2RequestReturnResponse makes a generic WSO2 request - returns err, http response, and response body
// toReturn should be a pointer
func MakeWSO2RequestReturnResponse(method, url string, body interface{}, toReturn interface{}) (*nerr.E, *http.Response, string) {
	err, response, responseBody := MakeWSO2RequestWithHeadersReturnResponse(method, url, body, toReturn, nil)
	return err, response, responseBody
}

// MakeWSO2RequestWithHeaders makes a generic WSO2 request with headers
// toReturn should be a pointer
func MakeWSO2RequestWithHeaders(method, url string, body interface{}, toReturn interface{}, headers map[string]string) *nerr.E {
	err, _, _ := MakeWSO2RequestWithHeadersReturnResponse(method, url, body, toReturn, headers)
	return err
}

// MakeWSO2RequestWithHeadersReturnResponse makes a generic WSO2 request with headers - returns err, http response, and response body
// toReturn should be a pointer
func MakeWSO2RequestWithHeadersReturnResponse(method, requestUrl string, body interface{}, toReturn interface{}, headers map[string]string) (*nerr.E, *http.Response, string) {
	log.L.Debugf("Making %v request against %v at %v", method, requestUrl, time.Now())

	key, er := GetAccessKey()
	if er != nil {
		return er.Addf("Couldn't make WSO2 request"), nil, ""
	}

	//attach key
	var b []byte
	var ok bool
	var err error
	hasRetried := false

	if body != nil {
		if b, ok = body.([]byte); !ok {
			b, err = json.Marshal(body)
			if err != nil {
				return nerr.Translate(err).Addf("Couldn't marhsal request"), nil, ""
			}
			log.L.Debugf("Sending %s to WSO2", b)
		}
	}

	for {
		req, err := http.NewRequest(method, requestUrl, bytes.NewBuffer(b))
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't build WSO2 request"), nil, ""
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", key))

		for k, v := range headers {
			log.L.Debugf("Setting header %v to %v", k, v)
			req.Header.Set(k, v)
		}

		c := http.Client{
			Timeout: 20 * time.Second, //I wish we could make this shorter... but alas.
		}

		resp, err := c.Do(req)
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return nerr.Translate(err).Addf("request timed out or was cancelled"), resp, ""
			}
		case *url.Error:
			if err, ok := err.Err.(net.Error); ok && err.Timeout() {
				return nerr.Translate(err).Addf("request timed out or was cancelled"), resp, ""
			}
		}

		if err != nil {
			return nerr.Translate(err).Addf("Couldn't make WSO2 request"), resp, ""
		}

		defer resp.Body.Close()

		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't read response body"), resp, ""
		}

		responseBody := string(rb)

		if resp.StatusCode/100 != 2 {
			if resp.StatusCode == 400 && len(rb) == 0 && !hasRetried {
				//if we get a 400 and a blank body and we haven't retried, then just try again
				log.L.Debugf("Retrying WSO2 request")
				hasRetried = true
				continue
			}

			return nerr.Create(fmt.Sprintf("response code %v: %s", resp.StatusCode, rb), "request-error"), resp, responseBody
		}

		err = json.Unmarshal(rb, toReturn)
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't unmarshal response %s", "unmarshal error"), resp, responseBody
		}

		return nil, resp, responseBody
	}
}
