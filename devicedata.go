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
import "encoding/hex"

// Socket data
type socketData struct {
	MsgType string      `json:"type"`
	Message string      `json:"message"`
	Data    DataMessage `json:"data"`
}

// DataMessage contains data from devices
type DataMessage struct {
	DeviceAddress  string  `json:"devAddr"`
	Timestamp      int64   `json:"timestamp"`
	StringData     string  `json:"data"`
	ApplicationEUI string  `json:"appEUI"`
	DeviceEUI      string  `json:"deviceEUI"`
	RSSI           int32   `json:"rssi"`
	SNR            float32 `json:"snr"`
	Frequency      float32 `json:"frequency"`
	GatewayEUI     string  `json:"gatewayEUI"`
	DataRate       string  `json:"dataRate"`
}

// Data returns the bytes sent by the device. If the bytes can't be parsed
// nil will be returned
func (d *DataMessage) Data() []byte {
	buf, _ := hex.DecodeString(d.StringData)
	return buf
}
