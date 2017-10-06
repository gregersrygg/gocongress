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
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
)

// DeviceType is used to identify the device type
type DeviceType string

const (
	// OTAA (Over-The-Air-Activation) devices use a join procedure to negotiate
	// new session keys and require just one application key
	OTAA = DeviceType("OTAA")
	// ABP (Activation By Personalization) devices have pre-provisioned session
	// keys and device addresses
	ABP = "ABP"
)

// Device is the type that represents your actual devices. A device must be
// associated with an application at all times. A device cannot be associated
// with more than one application at a time.
type Device struct {
	EUI                   string `json:"deviceEUI"`
	DeviceAddress         string `json:"devAddr"`
	ApplicationKey        string `json:"appKey"`
	ApplicationSessionKey string `json:"appSKey"`
	NetworkSessionKey     string `json:"nwkSKey"`
	FrameCounterUp        uint16 `json:"fCntUp"`
	FrameCounterDown      uint16 `json:"fCntDn"`
	RelaxedCounter        bool   `json:"relaxedCounter"`
	DeviceType            string `json:"deviceType"`
	KeyWarning            bool   `json:"keyWarning"`
	tagResource
	client *CongressClient
	app    *Application
}

// DownstreamMessage are messages sent to the devices
type DownstreamMessage struct {
	StringData  string `json:"data"`
	Port        uint8  `json:"port"`
	Ack         bool   `json:"ack"`
	SentTime    int64  `json:"sentTime"`
	CreatedTime int64  `json:"createdTime"`
	AckTime     int64  `json:"ackTime"`
	State       string `json:"state"`
}

// Data returns the bytes to be sent to the device
func (d *DownstreamMessage) Data() []byte {
	data, err := hex.DecodeString(d.StringData)
	if err != nil {
		return make([]byte, 0)
	}
	return data
}

var r *regexp.Regexp

func init() {
	var err error
	r, err = regexp.Compile("^[A-Za-z0-:_\\-+@\\ ,.=]*$")
	if err != nil {
		panic(fmt.Sprintf("I can't compile the string regexp: %v", err))
	}
}

type gwList struct {
	Gws []Gateway `json:"gateways"`
}

// UpstreamMessage is a message sent by the device to the backend
type UpstreamMessage struct {
	DeviceAddress string  `json:"devAddr"`
	Timestamp     int64   `json:"timestamp"`
	StringData    string  `json:"data"`
	AppEUI        string  `json:"appEUI"`
	DeviceEUI     string  `json:"deviceEUI"`
	RSSI          int32   `json:"rssi"`
	SNR           float32 `json:"snr"`
	Frequency     float32 `json:"frequency"`
	GatewayEUI    string  `json:"gatewayEUI"`
	DataRate      string  `json:"dataRate"`
}

// Data returns the bytes sent by the device
func (u *UpstreamMessage) Data() []byte {
	data, err := hex.DecodeString(u.StringData)
	if err != nil {
		return make([]byte, 0)
	}
	return data
}

// Update updates the device in the Congress backend. The updated device is returned.
func (device *Device) Update() (*Device, error) {
	ret, err := device.client.genericMutation(http.MethodPut, fmt.Sprintf("/applications/%s/devices/%s", device.app.EUI, device.EUI), device)
	if ret == nil {
		return nil, err
	}
	return ret.(*Device), err
}

// Delete removes the device from Congress
func (device *Device) Delete() error {
	return device.client.genericDelete(fmt.Sprintf("/applications/%s/devices/%s", device.app.EUI, device.EUI))
}

// EnqueueMessage enqueues a new downstream message to a device. The message will be sent the next
// time the device sends a packet upstream
func (device *Device) EnqueueMessage(data []byte, port uint8, ack bool) (*DownstreamMessage, error) {
	if port < 1 || port > 224 {
		return nil, ErrInvalidPort
	}
	newMsg := &DownstreamMessage{hex.EncodeToString(data), port, ack, 0, 0, 0, ""}
	ret, err := device.client.genericMutation(http.MethodPost, fmt.Sprintf("/applications/%s/devices/%s/message", device.app.EUI, device.EUI), newMsg)
	if ret == nil {
		return nil, err
	}
	return ret.(*DownstreamMessage), err
}

// GetQueuedMessage retrieves the currently queued downstream message
func (device *Device) GetQueuedMessage() (*DownstreamMessage, error) {
	msg := &DownstreamMessage{}
	ret, err := device.client.genericGet(fmt.Sprintf("/applications/%s/devices/%s/message", device.app.EUI, device.EUI), msg)
	if ret == nil {
		return nil, err
	}
	return ret.(*DownstreamMessage), err
}

// ClearEnqueuedMessage removes the enqueued downstream message
func (device *Device) ClearEnqueuedMessage() error {
	return device.client.genericDelete(fmt.Sprintf("/applications/%s/devices/%s/message", device.app.EUI, device.EUI))
}

// Messages returns the number of upstream messages sent from the device
func (device *Device) Messages(limit int) ([]UpstreamMessage, error) {

	type msgList struct {
		Msgs []UpstreamMessage `json:"messages"`
	}

	ret, err := device.client.genericGet(fmt.Sprintf("/applications/%s/devices/%s/data?limit=%d", device.app.EUI, device.EUI, limit), &msgList{})
	if err != nil {
		return nil, err
	}
	return ret.(*msgList).Msgs, err
}
