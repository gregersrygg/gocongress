package gocongress

/*
**   Copyright 2017 Telenor Digital AS
**
**  Licensed under the Apache License, Version 2.0 (the "License");
**  you may not use this file except in compliance with the License.
**  You may obtain a copy of the License at
**
**      http://www.apache.org/licenses/LICENSE-2.0
**
**  Unless required by applicable law or agreed to in writing, software
**  distributed under the License is distributed on an "AS IS" BASIS,
**  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
**  See the License for the specific language governing permissions and
**  limitations under the License.
 */

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
)

const (
	// DefaultEndpoint is the default endpoint used by Congress. It can be
	// overridden by setting the environment variable `CONGRESS_API_ENDPOINT`
	DefaultEndpoint = "https://api.lora.telenor.io"

	// EnvironmentToken is the default value for tokens; ie use environment variable
	// to retrieve API tokens.
	EnvironmentToken = ""

	tokenHeader = "X-API-Token"
)

// CongressClient is the client interface you use to interact with Congress.
type CongressClient struct {
	Token    string
	Endpoint string
	client   http.Client
}

// Create a new (default) request; set the content type and encode the entity
// into the request body if it is set.
func (c *CongressClient) newRequest(path string, entity interface{}) (*http.Request, error) {
	body := bytes.NewBufferString("")
	if entity != nil {
		if err := json.NewEncoder(body).Encode(entity); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(http.MethodGet, c.Endpoint+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(tokenHeader, c.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Ping performs a simple request to the root resource of the Congress server.
func (c *CongressClient) Ping() (*CongressClient, error) {
	_, err := c.genericGet("/", nil)
	return c, err
}

// Perform a generic GET request
func (c *CongressClient) genericGet(path string, entity interface{}) (interface{}, error) {
	req, err := c.newRequest(path, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, entity)
}

// Perform request, check errors and decode JSON response
func (c *CongressClient) doRequest(req *http.Request, entity interface{}) (interface{}, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := responseToError(resp); err != nil {
		return nil, err
	}
	if entity != nil {
		if err := json.NewDecoder(resp.Body).Decode(entity); err != nil {
			return nil, err
		}
	}
	return entity, nil
}

// Do a generic PUT or POST request with JSON in request and response body
func (c *CongressClient) genericMutation(method string, path string, entity interface{}) (interface{}, error) {
	req, err := c.newRequest(path, entity)
	if err != nil {
		return nil, err
	}
	req.Method = method
	return c.doRequest(req, entity)
}

// Do a generic DELETE - ie no content in request or response body.
func (c *CongressClient) genericDelete(path string) error {
	req, err := c.newRequest(path, nil)
	if err != nil {
		return err
	}
	req.Method = http.MethodDelete
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	return responseToError(resp)
}

// NewCongressClient creates a new CongressClient instance. If the token string
// is empty it will use the environment variable `CONGRESS_API_TOKEN`
func NewCongressClient(apiToken string) (*CongressClient, error) {
	ep := os.Getenv("CONGRESS_API_ENDPOINT")
	if ep == "" {
		ep = DefaultEndpoint
	}
	if apiToken == "" {
		apiToken = os.Getenv("CONGRESS_API_TOKEN")
	}
	client := &CongressClient{apiToken, ep, http.Client{}}

	return client.Ping()
}

// NewApplication creates a new application instance.
func (c *CongressClient) NewApplication() (*Application, error) {
	app := &Application{"", newTags(), c}
	ret, err := c.genericMutation(http.MethodPost, "/applications", app)
	if err != nil {
		return nil, err
	}
	return ret.(*Application), nil
}

// Applications return the list of your applications in Congress.
func (c *CongressClient) Applications() ([]Application, error) {
	type appList struct {
		Apps []Application `json:"applications"`
	}

	list, err := c.genericGet("/applications", &appList{})
	if err != nil {
		return nil, err
	}
	return list.(*appList).Apps, nil
}

// NewGateway creates a new gateway in Congress.
func (c *CongressClient) NewGateway(eui string, ip net.IP, strict bool, position *Position) (*Gateway, error) {
	gw := &Gateway{"", "", true, 0, 0, 0, newTags(), c}
	gw.EUI = eui
	gw.IP = ip.String()
	gw.StrictIP = strict
	if position != nil {
		gw.Altitude = position.Altitude
		gw.Latitude = position.Latitude
		gw.Longitude = position.Longitude
	}
	ret, err := c.genericMutation(http.MethodPost, "/gateways", gw)
	if err != nil {
		return nil, err
	}
	return ret.(*Gateway), nil
}

// Gateways return the list of your gateways in Congress.
func (c *CongressClient) Gateways() ([]Gateway, error) {
	list, err := c.genericGet("/gateways", &gwList{})
	if err != nil {
		return nil, err
	}
	return list.(*gwList).Gws, nil

}
