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
	"crypto/rand"
	"fmt"
	"net"
	"testing"
)

func randomEUI() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x-%02x-%02x",
		randomBytes[0], randomBytes[1], randomBytes[2], randomBytes[3],
		randomBytes[4], randomBytes[5], randomBytes[6], randomBytes[7])
}

func TestGateway(t *testing.T) {
	client, err := NewCongressClient(EnvironmentToken)
	if err != nil {
		t.Fatalf("Got error creating client: %v", err)
	}

	gw, err := client.NewGateway(randomEUI(), net.ParseIP("127.0.0.1"), false, nil)
	if err != nil {
		t.Fatalf("Got error creating new gateway: %v", err)
	}

	if !gw.SetTag("name", "REST Test Gateway") {
		t.Fatalf("Couldn't set gateway tag")
	}

	// Update it
	if _, err := gw.Update(); err != nil {
		t.Fatalf("Couldn't update gateway: %v", err)
	}

	// Retrieve list of gateways. It should be somewhere in the list
	gwList, err := client.Gateways()
	if err != nil {
		t.Fatalf("Got error retrieving list of gateways: %v", err)
	}
	found := false
	for _, v := range gwList {
		if v.EUI == gw.EUI {
			found = true
			if gw.GetTag("name") != v.GetTag("name") {
				t.Fatalf("Couldn't find name tag. Got %s but expected %s", v.GetTag("name"), gw.GetTag("name"))
			}
		}
	}
	if !found {
		t.Fatal("Couldn't locate the gateway in the list")
	}

	if err := gw.Delete(); err != nil {
		t.Fatalf("Got error removing gateway: %v", err)
	}
}
