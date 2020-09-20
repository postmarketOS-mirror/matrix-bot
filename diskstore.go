// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
	"os"
	"path/filepath"
)

// DiskStore implements the Storer interface.
//
// Data is persisted in a json file with the provided filename.
type DiskStore struct {
	File      string `json:"-"`
	NextBatch string `json:"next_batch"`
	FilterID  string `json:"filter_id"`
}

func (config *DiskStore) Load() {
	config.load(config.File, &config)
}

func (config *DiskStore) Save() {
	config.save(config.File, &config)
}

func (config *DiskStore) load(file string, target interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Failed to read", file)
		panic(err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		fmt.Println("Failed to parse", file)
		panic(err)
	}
}

func (config *DiskStore) save(file string, source interface{}) {
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		fmt.Println("Failed to create", dir)
		panic(err)
	}

	data, err := json.Marshal(source)
	if err != nil {
		fmt.Println("Failed to marshal data for", file)
		panic(err)
	}

	err = ioutil.WriteFile(file, data, 0600)
	if err != nil {
		fmt.Println("Failed to write to", file)
		panic(err)
	}
}

// Storer interface

func (config *DiskStore) SaveFilterID(_ id.UserID, filterID string) {
	config.FilterID = filterID
	config.Save()
}

func (config *DiskStore) LoadFilterID(_ id.UserID) string {
	return config.FilterID
}

func (config *DiskStore) SaveNextBatch(_ id.UserID, nextBatch string) {
	config.NextBatch = nextBatch
	config.Save()
}

func (config *DiskStore) LoadNextBatch(_ id.UserID) string {
	return config.NextBatch
}

func (config *DiskStore) SaveRoom(_ *mautrix.Room) {
	panic("SaveRoom is not implemented")
}

func (config *DiskStore) LoadRoom(_ id.RoomID) *mautrix.Room {
	panic("LoadRoom is not implemented")
}

func NewDiskStore(file string) *DiskStore {
	return &DiskStore{
		File: file,
	}
}
