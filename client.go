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
	"fmt"
	"net"
	"net/http"
)

const (
	// DefaultAddr is the default address of Congress.
	DefaultAddr = "https://api.lora.telenor.io"

	tokenHeader = "X-API-Token"
)

// CongressClient is the client interface you use to interact with Congress.
type CongressClient struct {
	Addr   string
	Token  string
	client http.Client
}

// NewCongressClient creates a new CongressClient.
func NewCongressClient(token string) (*CongressClient, error) {
	return NewCongressClientWithAddr(DefaultAddr, token)
}

// NewCongressClientWithAddr creates a new CongressClient that talks to the supplied address.
// It is only useful for internal testing; all clients should use NewCongressClient.
func NewCongressClientWithAddr(addr, token string) (*CongressClient, error) {
	c := &CongressClient{
		Addr:  addr,
		Token: token,
	}
	return c, c.Ping()
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
	req, err := http.NewRequest(http.MethodGet, c.Addr+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(tokenHeader, c.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Ping performs a simple request to the root resource of the Congress server.
func (c *CongressClient) Ping() error {
	_, err := c.genericGet("/", nil)
	return err
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

// NewApplication creates a new application.
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

// GetApplication retrieves an application from Congress.
func (c *CongressClient) GetApplication(eui string) (*Application, error) {
	app := &Application{"", newTags(), c}
	existingApp, err := c.genericGet(fmt.Sprintf("/applications/%s", eui), app)
	if err != nil {
		return nil, err
	}
	return existingApp.(*Application), nil
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
