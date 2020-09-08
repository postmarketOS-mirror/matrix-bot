// pmos-bot - A bot for the postmarketOS Matrix channels
// Copyright (C) 2017 Tulir Asokan
// Copyright (C) 2018-2019 Luca Weiss
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"os"
	"regexp"
	"strings"
)

var homeserver = flag.String("homeserver", "https://matrix.org", "Matrix homeserver")
var username = flag.String("username", "", "Matrix username localpart")
var password = flag.String("password", "", "Matrix password")
var deviceId = flag.String("deviceid", "", "Matrix device id (optional)")
var diskStorePath = flag.String("stateStoragePath", "", "Path to a .json file where state information is stored")

var shortcutMapRegex = buildShortcutMapRegex()

func main() {
	flag.Parse()
	if *username == "" || *password == "" || *diskStorePath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("Logging to", *homeserver, "as", *username)
	client, err := mautrix.NewClient(*homeserver, "", "")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	diskStore := NewDiskStore(*diskStorePath)
	diskStore.Load()
	client.Store = diskStore

	_, err = client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: "m.id.user", User: *username},
		Password:         *password,
		DeviceID:         id.DeviceID(*deviceId),
		StoreCredentials: true,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Println("Login successful")

	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		content, err := handleEvent(evt)
		if err != nil {
			fmt.Println(err)
			return
		}

		// No content - nothing to do
		if content == nil {
			return
		}

		_, err = client.SendMessageEvent(evt.RoomID, event.EventMessage, content)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	err = client.Sync()
	if err != nil {
		fmt.Println(err)
	}
}

func handleEvent(evt *event.Event) (*event.MessageEventContent, error) {
	senderName, _, err := evt.Sender.Parse()
	if err != nil {
		return nil, err
	}
	if senderName != *username {
		var body string
		// Use FormattedBody is available, as it will contain quote information that we want to remove
		if len(evt.Content.AsMessage().FormattedBody) != 0 {
			msg := evt.Content.AsMessage()
			msg.RemoveReplyFallback()
			body = msg.FormattedBody
		} else {
			body = evt.Content.AsMessage().Body
		}
		//fmt.Printf("DBG <%[1]s> %[4]s (%[2]s/%[3]s)\n", evt.Sender, evt.Type.String(), evt.ID, body)
		matches := shortcutMapRegex.FindAllStringSubmatch(body, -1)
		if matches != nil {
			var buffer bytes.Buffer
			for _, match := range matches {
				//fmt.Println(match[1] + match[2] + " matched!")
				fmt.Printf("<%[1]s> %[4]s (%[2]s/%[3]s)\n", evt.Sender, evt.Type.String(), evt.ID, body)
				buffer.WriteString(shortcutMap[strings.ToLower(match[1])] + match[2] + " ")
			}
			return &event.MessageEventContent{MsgType: event.MsgText, Body: strings.TrimSuffix(buffer.String(), " ")}, nil
		}
	}
	return nil, nil
}

func buildShortcutMapRegex() *regexp.Regexp {
	keys := make([]string, len(shortcutMap))
	i := 0
	for k := range shortcutMap {
		keys[i] = k
		i++
	}

	regex := "(?i)(" + strings.Join(keys, "|") + ")(\\d+)"
	return regexp.MustCompile(regex)
}
