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
	"testing"
)

func TestDevices(t *testing.T) {
	client, err := NewCongressClient(EnvironmentToken)
	if err != nil {
		t.Fatalf("Couldn't create Congress client: %v", err)
	}

	app, _ := client.NewApplication()
	otaa, err := app.NewDevice(OTAA)

	if err != nil {
		t.Fatalf("Couldn't create OTAA device: %v", err)
	}

	otaa.SetTag("name", "REST Client Test OTAA")
	otaa.RelaxedCounter = true
	otaa.FrameCounterDown = 0
	otaa.FrameCounterUp = 0

	updatedDevice, err := otaa.Update()
	if err != nil {
		t.Fatalf("Couldn't update device: %v", err)
	}

	if updatedDevice.GetTag("name") != otaa.GetTag("name") {
		t.Fatalf("Tag is different on updated device. Expected %s but got %s", otaa.GetTag("name"), updatedDevice.GetTag("name"))
	}

	_, err = updatedDevice.Messages(60)
	if err != nil {
		t.Fatalf("Got error retrieving messages: %v", err)
	}

	if err := updatedDevice.Delete(); err != nil {
		t.Fatalf("Couldn't delete device: %v", err)
	}

	if err := otaa.Delete(); ErrorStatusCode(err) != http.StatusNotFound {
		t.Fatalf("Expected not found on deleted device but got %v", err)
	}

	abp, err := app.NewDevice(ABP)
	if err != nil {
		t.Fatalf("Couldn't create ABP device: %v", err)
	}
	if err := abp.Delete(); err != nil {
		t.Fatalf("Got error removing ABP device: %v", err)
	}
	if err := abp.Delete(); ErrorStatusCode(err) != http.StatusNotFound {
		t.Fatalf("Expected ErrNotFound but got %v", err)
	}
}

func TestDownstreamMessages(t *testing.T) {
	client, err := NewCongressClient(EnvironmentToken)
	if err != nil {
		t.Fatalf("Couldn't create Congress client: %v", err)
	}

	app, _ := client.NewApplication()
	device, _ := app.NewDevice(OTAA)

	msg, err := device.EnqueueMessage([]byte{1, 2, 3, 4, 5, 6, 7}, 1, false)
	if err != nil {
		t.Fatalf("Got error queuing downstream message: %v", err)
	}

	qMsg, err := device.GetQueuedMessage()
	if err != nil {
		t.Fatalf("Couldn't retrieve enqueued message: %v", err)
	}
	if qMsg.StringData != msg.StringData {
		t.Fatalf("Data is different.")
	}
	if err := device.ClearEnqueuedMessage(); err != nil {
		t.Fatalf("Couldn't remove downstream message: %v", err)
	}
}

// Create a lot of devices, then update them individually
func TestMultipleDevices(t *testing.T) {
	client, _ := NewCongressClient(EnvironmentToken)
	app, _ := client.NewApplication()

	devices := make([]*Device, 10)
	var err error
	for i := 0; i < len(devices); i++ {
		if devices[i], err = app.NewDevice(ABP); err != nil {
			t.Errorf("Got error creating device: %v", err)
		}
	}

	// Update each device separately
	for i := 0; i < len(devices); i++ {
		devices[i].DeviceAddress = fmt.Sprintf("%08x", i)
		if devices[i], err = devices[i].Update(); err != nil {
			t.Errorf("Got error updating device: %v", err)
		}
	}

	for i := 0; i < len(devices); i++ {
		if devices[i].DeviceAddress != fmt.Sprintf("%08x", i) {
			t.Errorf("Did not get the expected DevAddr for device %d: Got %s expected %s", i, devices[i].DeviceAddress, fmt.Sprintf("%08x", i))
		}
	}

	for i := 0; i < len(devices); i++ {
		if err := devices[i].Delete(); err != nil {
			t.Errorf("Got error deleting device %d: %v", i, err)
		}
	}

	app.Delete()
}
