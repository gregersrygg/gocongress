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
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	addr  = flag.String("addr", DefaultAddr, "congress API addr")
	token = flag.String("api-token", "", "congress API token")
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *token == "" {
		if _, err := NewCongressClientWithAddr(*addr, ""); err != nil {
			fmt.Println("Error creating client:", err)
			fmt.Println("You might need to set the token flag when running the tests against the address", *addr)
			fmt.Println("That is, run `go test -args -token <your-api-token>.")
			os.Exit(1)
		}
	}

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	// Ping is implicit
	client, err := NewCongressClientWithAddr(*addr, *token)
	if err != nil {
		t.Fatalf("Got error creating client: %v", err)
	}

	if err := client.Ping(); err != nil {
		t.Fatalf("Got error calling ping(): %v", err)
	}
}
