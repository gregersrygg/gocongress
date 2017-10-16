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
	"reflect"
	"testing"
)

func TestDeviceData(t *testing.T) {
	d := DataMessage{StringData: "BEEFBABE"}

	if !reflect.DeepEqual(d.Data(), []byte{0xBE, 0xEF, 0xBA, 0xBE}) {
		t.Fatal("Couldn't parse bytes")
	}

	d.StringData = "invalid hex char"
	if d.Data() != nil {
		t.Fatal("Couldn't parse bytes")
	}
}
