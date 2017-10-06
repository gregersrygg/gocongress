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

import "testing"

func TestTags(t *testing.T) {
	tr := newTags()
	if !tr.SetTag("name", "Tag test") {
		t.Fatal("Couldn't set tag")
	}
	if tr.GetTag("name") != "Tag test" || tr.GetTag("nAmE") != "Tag test" {
		t.Fatal("Name tag isn't the same")
	}

	if tr.SetTag("name", "alert('Hello world');") {
		t.Fatal("Could set illegal characters in tag value")
	}

	if tr.SetTag("alert('Hello');", "test") {
		t.Fatal("Could set illegal characters in tag name")
	}
}
