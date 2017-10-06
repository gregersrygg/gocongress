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
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("CONGRESS_API_TOKEN") == "" {
		fmt.Println("The CONGRESS_API_TOKEN environment variable must be set to run the tests")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func xTestPing(t *testing.T) {
	// Ping is implicit
	if _, err := NewCongressClient("invalid"); err == nil {
		t.Fatal("Expected error when using invalid token but the returned error was nil")
	}

	client, err := NewCongressClient(EnvironmentToken)
	if err != nil {
		t.Fatalf("Got error creating client: %v", err)
	}

	if _, err := client.Ping(); err != nil {
		t.Fatalf("Got error calling ping(): %v", err)
	}
}
