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
)

// Gateway is an entity representing gateways in Congress. Gateways are the
// main link between your devices and your backend software. They forward the
// radio packets from the devices. A single gateway will forward packets
// to any and all applications in the backend.
type Gateway struct {
	EUI       string  `json:"gatewayEUI,omitempty"`
	IP        string  `json:"ip,omitempty"`
	StrictIP  bool    `json:"strictIP"`
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
	Altitude  float32 `json:"altitude,omitempty"`
	tagResource
	client *CongressClient
}

// Position represents a geographical position with latitude, longitude and altitude
type Position struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
}

// Update updates the gateway in the Congress backend. The updated gateway is returned.
func (gw *Gateway) Update() (*Gateway, error) {
	ret, err := gw.client.genericMutation(http.MethodPut, fmt.Sprintf("/gateways/%s", gw.EUI), gw)
	if ret == nil {
		return nil, err
	}
	return ret.(*Gateway), err
}

// Delete removes the gateway from the Congress backend.
func (gw *Gateway) Delete() error {
	return gw.client.genericDelete(fmt.Sprintf("/gateways/%s", gw.EUI))
}
