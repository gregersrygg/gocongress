// Package gocongress is a client library for the Congress LoRa backend.
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
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

// Application is the type that represent an application in Congress. LoRa
// applications are just a way to group related devices into groups.
type Application struct {
	// EUI is the application's EUI
	EUI string `json:"applicationEUI,omitempty"`
	tagResource
	client *CongressClient
}

// AppOutput is an application output
type AppOutput struct {
	EUI    string                 `json:"eui,omitempty"`
	AppEUI string                 `json:"appEUI,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
	Log    []OutputLog            `json:"logs,omitempty"`
	Status string                 `json:"status,omitempty"`
	app    *Application
	client *CongressClient
}

// OutputLog is the log from the output
type OutputLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// OutputConfig is a generic output configuration
type OutputConfig interface {
	Config() map[string]interface{}
}

// MQTTConfig is a configuration struct for MQTT
type MQTTConfig struct {
	Endpoint         string `json:"endpoint"`
	Port             int    `json:"port"`
	TLS              bool   `json:"tls,omitempty"`
	CertificateCheck bool   `json:"certCheck,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	ClientID         string `json:"clientid,omitempty"`
	TopicName        string `json:"topicName,omitempty"`
}

// Config returns the configuration fields as a map
func (m *MQTTConfig) Config() map[string]interface{} {
	return map[string]interface{}{
		"type":      "mqtt",
		"endpoint":  m.Endpoint,
		"port":      m.Port,
		"tls":       m.TLS,
		"certCheck": m.CertificateCheck,
		"username":  m.Username,
		"password":  m.Password,
		"clientid":  m.ClientID,
		"topicName": m.TopicName,
	}
}

func (m *MQTTConfig) string(t map[string]interface{}, key string, def string) string {
	val, ok := t[key].(string)
	if !ok {
		return def
	}
	return val
}

func (m *MQTTConfig) bool(t map[string]interface{}, key string, def bool) bool {
	val, ok := t[key].(bool)
	if !ok {
		return def
	}
	return val
}

func (m *MQTTConfig) int(t map[string]interface{}, key string, def int) int {
	val, ok := t[key].(float64)
	if !ok {
		return def
	}
	return int(val)
}

// ReadFromMap applies the configuration in the map to the fields
func (m *MQTTConfig) ReadFromMap(vals map[string]interface{}) {
	m.Endpoint = m.string(vals, "endpoint", "")
	m.Port = m.int(vals, "port", 1883)
	m.TLS = m.bool(vals, "tls", false)
	m.CertificateCheck = m.bool(vals, "certCheck", true)
	m.Username = m.string(vals, "username", "")
	m.Password = m.string(vals, "password", "")
	m.ClientID = m.string(vals, "clietnid", "")
	m.TopicName = m.string(vals, "topicName", "")
}

// Update updates the application in the Congress backend. The updated application
// is returned.
func (app *Application) Update() (*Application, error) {
	res, err := app.client.genericMutation(http.MethodPut, fmt.Sprintf("/applications/%s", app.EUI), app)
	if res == nil {
		return nil, err
	}
	return res.(*Application), err
}

// Delete removes the application from Congress
func (app *Application) Delete() error {
	return app.client.genericDelete(fmt.Sprintf("/applications/%s", app.EUI))
}

// NewDevice creates a new OTAA (Over-The-Air-Activated) device in Congress.
// The AppKey and EUI are automatically generated by the Congress backend.
func (app *Application) NewDevice(dt DeviceType) (*Device, error) {
	device := &Device{"", "", "", "", "", 0, 0, false, "", false, newTags(), app.client, app}
	if dt == OTAA {
		device.DeviceType = "OTAA"
	} else {
		device.DeviceType = "ABP"
	}
	ret, err := app.client.genericMutation(http.MethodPost, fmt.Sprintf("/applications/%s/devices", app.EUI), device)
	if err != nil {
		return nil, err
	}
	return ret.(*Device), nil
}

// Outputs returns the list of configured outputs
func (app *Application) Outputs() ([]AppOutput, error) {
	type outputList struct {
		Outputs []AppOutput `json:"outputs"`
	}
	list, err := app.client.genericGet(fmt.Sprintf("/applications/%s/outputs", app.EUI), &outputList{})
	if err != nil {
		return nil, err
	}
	return list.(*outputList).Outputs, nil
}

// Devices returns the device list for the application
func (app *Application) Devices() ([]Device, error) {
	type deviceList struct {
		Devices []Device `json:"devices"`
	}
	list, err := app.client.genericGet(fmt.Sprintf("/applications/%s/devices", app.EUI), &deviceList{})
	if err != nil {
		return nil, err
	}
	return list.(*deviceList).Devices, nil
}

// NewOutput creates a new application output
func (app *Application) NewOutput(config OutputConfig) (*AppOutput, error) {
	output := &AppOutput{Config: config.Config(), app: app, client: app.client}
	ret, err := app.client.genericMutation(http.MethodPost, fmt.Sprintf("/applications/%s/outputs", app.EUI), output)
	if err != nil {
		return nil, err
	}
	return ret.(*AppOutput), nil
}

// Update updates the application output
func (output *AppOutput) Update() (*AppOutput, error) {
	res, err := output.client.genericMutation(http.MethodPut, fmt.Sprintf("/applications/%s/outputs/%s", output.app.EUI, output.EUI), output)
	if res == nil {
		return nil, err
	}
	return res.(*AppOutput), err
}

// Delete removes the application output
func (output *AppOutput) Delete() error {
	return output.client.genericDelete(fmt.Sprintf("/applications/%s/outputs/%s", output.app.EUI, output.EUI))
}

// DataErrorMessage are error messages generated by the data stream.
type DataErrorMessage string

// DataStream returns a channel with device data using the appliction's web
// socket. If there's an error reading the web socket the channel will be closed.
// Error messages are sent on the error channel that is returned.
func (app *Application) DataStream() (chan DataMessage, chan DataErrorMessage, error) {

	congressURL, err := url.Parse(app.client.Addr)
	if err != nil {
		return nil, nil, err
	}

	wscfg, err := websocket.NewConfig(fmt.Sprintf("wss://%s/applications/%s/stream", congressURL.Host, app.EUI), "http://example.com")
	if err != nil {
		return nil, nil, err
	}
	wscfg.Header.Set(tokenHeader, app.client.Token)

	ws, err := websocket.DialConfig(wscfg)
	if err != nil {
		return nil, nil, err
	}

	ret := make(chan DataMessage)
	errors := make(chan DataErrorMessage)
	go func() {
		defer ws.Close()
		defer close(ret)
		defer close(errors)
		for {
			data := socketData{}
			err := websocket.JSON.Receive(ws, &data)
			if err != nil {
				errors <- DataErrorMessage(fmt.Sprintf("%v", err))
				return
			}
			switch data.MsgType {
			case "DeviceData":
				select {
				case ret <- data.Data:
				case <-time.After(400 * time.Millisecond):
					errors <- DataErrorMessage("Timed out writing to socket")
					return
				}
			case "Error":

				return
			default:
				// Ignore it
			}
		}
	}()
	return ret, errors, nil
}
