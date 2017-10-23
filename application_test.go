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
	"net/http"
	"testing"
	"time"
)

func TestApplications(t *testing.T) {
	client, err := NewCongressClientWithAddr(*addr, *token)
	if err != nil {
		t.Fatalf("Got error creating client: %v", err)
	}
	_, err = client.Applications()
	if err != nil {
		t.Fatalf("Got error retrieving app list: %v", err)
	}

	// Create a new application
	app, err := client.NewApplication()
	if err != nil {
		t.Fatalf("Got error retrieving application: %v", err)
	}
	// Update the application
	app.SetTag("name", "REST Library Test")

	updatedApp, err := app.Update()
	if err != nil {
		t.Fatalf("Got error updating application: %v", err)
	}
	if updatedApp.GetTag("name") != app.GetTag("name") {
		t.Fatal("Updated app does not contain the name tag")
	}
	// Retrieve the application list again. The updated application
	// should be somewhere in the returned list
	list, err := client.Applications()
	if err != nil {
		t.Fatalf("Couldn't retrieve applications: %v", err)
	}
	found := false
	for _, a := range list {
		if a.EUI == app.EUI {
			found = true
			if a.GetTag("name") != updatedApp.GetTag("name") {
				t.Fatalf("Name tag isn't matching. Got %s but expected %s", a.GetTag("name"), app.GetTag("name"))
			}
		}
	}
	if !found {
		t.Fatal("Couldn't locate application in list")
	}

	if _, err := client.GetApplication(app.EUI); err != nil {
		t.Fatal("Should be able to retrieve application")
	}

	// Delete the application. List shouldn't contain any applications
	if err := app.Delete(); err != nil {
		t.Fatalf("Couldn't delete application: %v", err)
	}
	if err := app.Delete(); ErrorStatusCode(err) != http.StatusNotFound {
		t.Fatalf("Expected ErrNotFound when deleting app for a second time but got: %v", err)
	}
}

func TestAppOutput(t *testing.T) {
	client, _ := NewCongressClientWithAddr(*addr, *token)
	app, _ := client.NewApplication()

	mqtt1 := &MQTTConfig{
		Endpoint:  "localhost",
		Port:      1883,
		TLS:       false,
		Username:  "john",
		Password:  "doe",
		ClientID:  "congress",
		TopicName: "testOutput",
	}
	mqtt2 := &MQTTConfig{
		Endpoint:  "localhost",
		Port:      1883,
		TLS:       false,
		Username:  "john",
		Password:  "doe",
		ClientID:  "congress",
		TopicName: "testOutput",
	}

	// Create two outputs

	op1, err := app.NewOutput(mqtt1)
	if err != nil {
		t.Fatalf("Couldn't create output: %v", err)
	}
	if op1 == nil {
		t.Fatal("No output returned")
	}

	op2, err := app.NewOutput(mqtt2)
	if err != nil {
		t.Fatalf("Couldn't create output: %v", err)
	}
	if op2 == nil {
		t.Fatal("No output returned")
	}

	// Retrieve the list. Should contain both
	opList, err := app.Outputs()
	if err != nil {
		t.Fatalf("Couldn't retrieve list of outputs: %v", err)
	}

	if len(opList) != 2 {
		t.Fatalf("Output list contains %d elements. Expected 2.", len(opList))
	}

	// Update config on 1 and 2, ensure they are updated
	mqtt1.Endpoint = "first"
	op1.Config = mqtt1.Config()
	newOp1, err := op1.Update()
	if err != nil {
		t.Fatalf("Got error updating output 1: %v", err)
	}
	if newOp1.Config["endpoint"] != "first" {
		t.Fatal("Endpoint didn't update")
	}

	mqtt2.Endpoint = "second"
	op2.Config = mqtt2.Config()
	newOp2, err := op2.Update()
	if err != nil {
		t.Fatalf("Got error updating output 2: %v", err)
	}
	if newOp2.Config["endpoint"] != "second" {
		t.Fatal("Endpoint for op2 didn't update")
	}

	if err := newOp1.Delete(); err != nil {
		t.Fatalf("Couldn't delete output 1 EUI=%s", newOp1.EUI)
	}
	if err := newOp2.Delete(); err != nil {
		t.Fatalf("Couldn't delete output 2 EUI=%s", newOp2.EUI)
	}

	if err := op1.Delete(); ErrorStatusCode(err) != http.StatusNotFound {
		t.Fatalf("Expected 404 not found when deleting output a second time but got %d", ErrorStatusCode(err))
	}

	app.Delete()
}

// Test the websocket output. There's no easy way to generate output
func TestWebsocketOutput(t *testing.T) {
	client, _ := NewCongressClientWithAddr(*addr, *token)
	app, _ := client.NewApplication()
	d1, _ := app.NewDevice(OTAA)
	d2, _ := app.NewDevice(ABP)

	ch, _, err := app.DataStream()
	if err != nil {
		t.Fatalf("Couldn't open data stream for application: %v", err)
	}

	select {
	case <-ch:
	case <-time.After(700 * time.Millisecond):
		// OK
	}
	d1.Delete()
	d2.Delete()

	app.Delete()
}
