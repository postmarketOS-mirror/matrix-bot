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
	"os"
	"regexp"
	"strings"
)

var homeserver = flag.String("homeserver", "https://matrix.org", "Matrix homeserver")
var username = flag.String("username", "", "Matrix username localpart")
var password = flag.String("password", "", "Matrix password")
var deviceId = flag.String("deviceid", "", "Matrix device id (optional)")

func main() {
	flag.Parse()
	if *username == "" || *password == "" {
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
	resp, err := client.Login(&mautrix.ReqLogin{
		Type:       "m.login.password",
		Identifier: mautrix.UserIdentifier{Type: "m.id.user", User: *username},
		Password:   *password,
		DeviceID:   *deviceId,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	client.SetCredentials(resp.UserID, resp.AccessToken)

	fmt.Println("Login successful")

	shortcutmap := map[string]string{
		"art#": "https://gitlab.com/postmarketOS/artwork/issues/",
		"art!": "https://gitlab.com/postmarketOS/artwork/merge_requests/",

		"bpo#": "https://gitlab.com/postmarketOS/build.postmarketos.org/issues/",
		"bpo!": "https://gitlab.com/postmarketOS/build.postmarketos.org/merge_requests/",

		"bot#": "https://gitlab.com/postmarketOS/matrix-bot/issues/",
		"bot!": "https://gitlab.com/postmarketOS/matrix-bot/merge_requests/",

		"chrg#": "https://gitlab.com/postmarketOS/charging-sdl/issues/",
		"chrg!": "https://gitlab.com/postmarketOS/charging-sdl/merge_requests/",

		"lnx#": "https://gitlab.com/postmarketOS/linux-postmarketos/issues/",
		"lnx!": "https://gitlab.com/postmarketOS/linux-postmarketos/merge_requests/",

		"mrh#": "https://gitlab.com/postmarketOS/mrhlpr/issues/",
		"mrh!": "https://gitlab.com/postmarketOS/mrhlpr/merge_requests/",

		"osk#": "https://gitlab.com/postmarketOS/osk-sdl/issues/",
		"osk!": "https://gitlab.com/postmarketOS/osk-sdl/merge_requests/",

		"pma#": "https://gitlab.com/postmarketOS/pmaports/issues/",
		"pma!": "https://gitlab.com/postmarketOS/pmaports/merge_requests/",

		"pmb#": "https://gitlab.com/postmarketOS/pmbootstrap/issues/",
		"pmb!": "https://gitlab.com/postmarketOS/pmbootstrap/merge_requests/",

		"org#": "https://gitlab.com/postmarketOS/postmarketos.org/issues/",
		"org!": "https://gitlab.com/postmarketOS/postmarketos.org/merge_requests/",

		"wiki#": "https://gitlab.com/postmarketOS/wiki/issues/",
	}
	shortcutmapregex := regexp.MustCompile("(?i)(art[#!]|bpo[#!]|bot[#!]|chrg[#!]|lnx[#!]|mrh[#!]|osk[#!]|pma[#!]|pmb[#!]|org[#!]|wiki#)(\\d+)")

	removequoteregex := regexp.MustCompile("<mx-reply>.*</mx-reply>")

	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(mautrix.EventMessage, func(evt *mautrix.Event) {
		if evt.Sender != *username {
			var body string
			// Use FormattedBody is available, as it will contain quote information that we want to remove
			if len(evt.Content.FormattedBody) != 0 {
				body = evt.Content.FormattedBody
				body = removequoteregex.ReplaceAllString(body, "")
			} else {
				body = evt.Content.Body
			}
			matches := shortcutmapregex.FindAllStringSubmatch(body, -1)
			if matches != nil {
				var buffer bytes.Buffer
				for _, match := range matches {
					fmt.Println(match[1] + match[2] + " matched!")
					fmt.Printf("<%[1]s> %[4]s (%[2]s/%[3]s)\n", evt.Sender, evt.Type.String(), evt.ID, body)
					buffer.WriteString(shortcutmap[strings.ToLower(match[1])] + match[2] + " ")
				}
				content := mautrix.Content{MsgType: mautrix.MsgText, Body: strings.TrimSuffix(buffer.String(), " ")}
				_, err := client.SendMessageEvent(evt.RoomID, mautrix.EventMessage, content)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})

	err = client.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
