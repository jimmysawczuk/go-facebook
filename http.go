package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPMethod string

const (
	Get  HTTPMethod = "GET"
	Post            = "POST"
	Put             = "PUT"
)

// Makes a standard API call to the Graph API.
func (f *Client) api(req_url string, method HTTPMethod, params map[string]interface{}) (map[string]interface{}, error) {

	if params == nil {
		params = make(map[string]interface{})
	}

	if _, exists := params["access_token"]; !exists && !f.accessToken.Empty() && f.accessToken.Valid() {
		params["access_token"] = f.accessToken.token
	}

	req_url = graph_endpoint + req_url + getQueryString(params)

	resp, err := http.Get(req_url)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		g_err, err := parseError(resp.StatusCode, buf)
		if err == nil {
			return nil, g_err
		} else {
			return nil, err
		}
	} else {

		map_result := make(map[string]interface{})
		err = json.Unmarshal(buf, &map_result)
		if err == nil {
			return map_result, err
		}

		values, err := url.ParseQuery(string(buf))
		if err == nil {
			for k, v := range values {
				if len(v) == 1 {
					map_result[k] = v[0]
				} else {
					map_result[k] = v
				}
			}

			return map_result, nil
		}

		return map_result, fmt.Errorf("couldn't parse response from graph")
	}
}

// A wrapper method for Api to make POST requests.
func (f *Client) Post(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Post, params)
}

// A wrapper method for Api to make GET requests.
func (f *Client) Get(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Get, params)
}

// A wrapper method for Api to make PUT requests.
func (f *Client) Put(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Put, params)
}

func getQueryString(params map[string]interface{}) string {
	values, _ := url.ParseQuery("")

	for key, value := range params {
		switch value.(type) {
		case string:
			values.Add(key, value.(string))
		case fmt.Stringer:
			values.Add(key, value.(fmt.Stringer).String())
		default:
			panic("Can't make a string!")
		}
	}

	result := values.Encode()

	if result != "" {
		result = "?" + result
	}

	return result
}
