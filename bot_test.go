// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"encoding/json"
	"io/ioutil"
	"maunium.net/go/mautrix/event"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func parseEventFromFile(filename string) event.Event {
	data, err := ioutil.ReadFile(filename)
	check(err)

	evt := event.Event{}
	err = json.Unmarshal(data, &evt)
	check(err)

	err = evt.Content.ParseRaw(evt.Type)
	check(err)

	return evt
}

func TestPlainMessage(t *testing.T) {
	evt := parseEventFromFile("./testdata/plain_message.json")

	content, err := handleEvent(&evt)
	check(err)
	if content != nil {
		t.Errorf("content should be nil")
	}
}

func TestMessageKeyword(t *testing.T) {
	var tests = []struct {
		filename string
		message  string
	}{
		{"./testdata/message_keyword.json", "https://gitlab.com/postmarketOS/pmaports/issues/123"},
		{"./testdata/message_keyword2.json", "https://gitlab.com/postmarketOS/pmbootstrap/merge_requests/456"},
	}

	for _, tt := range tests {
		evt := parseEventFromFile(tt.filename)

		content, err := handleEvent(&evt)
		check(err)

		assertEqual(t, content.MsgType, event.MsgNotice)
		assertEqual(t, content.Body, tt.message)
	}
}

func TestReply(t *testing.T) {
	evt := parseEventFromFile("./testdata/reply.json")

	content, err := handleEvent(&evt)
	check(err)
	if content != nil {
		t.Errorf("content should be nil")
	}
}
