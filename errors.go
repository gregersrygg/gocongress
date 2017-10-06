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
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// ErrNotImplemented is a placeholder for greater things to come
	ErrNotImplemented = &CongressError{Message: "Not implemented yet", StatusCode: http.StatusTeapot}
	// ErrInvalidPort is returned when the downstream port number is invalid
	ErrInvalidPort = &CongressError{Message: "Invalid port number", StatusCode: http.StatusBadRequest}
)

// CongressError contains the error messages emitted by Congress
type CongressError struct {
	Message    string
	StatusCode int
}

func (c *CongressError) Error() string {
	return fmt.Sprintf("%d: %s", c.StatusCode, c.Message)
}

// Create a new CongressError instance from a response.
func newCongressError(resp *http.Response) *CongressError {
	ret := &CongressError{StatusCode: resp.StatusCode}
	if buf, err := ioutil.ReadAll(resp.Body); err != nil {
		ret.Message = err.Error()
	} else {
		ret.Message = string(buf)
	}
	return ret
}

// Convert http response to error
func responseToError(response *http.Response) error {
	if response.StatusCode < 300 {
		return nil
	}
	return newCongressError(response)
}

// ErrorMessage returns the message part of the CongressError error
func ErrorMessage(err error) string {
	msg, ok := err.(*CongressError)
	if ok {
		return msg.Message
	}
	return ""
}

// ErrorStatusCode returns the HTTP status code of the CongressError error
func ErrorStatusCode(err error) int {
	msg, ok := err.(*CongressError)
	if ok {
		return msg.StatusCode
	}
	return 0
}
