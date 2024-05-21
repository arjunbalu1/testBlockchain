package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

// Define a struct to match the JSON structure of your values
type Entry struct {
	Key   string `json:"key,omitempty"`
	Value int    `json:"value,omitempty"`
	Ver   int    `json:"ver"`
	Valid bool   `json:"valid,omitempty"`
	Hash  string `json:"hash,omitempty"`
	Val   int    `json:"val,omitempty"`
}

func main() {
	db, err := leveldb.OpenFile("../leveldb_data", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var entry Entry
		if err := json.Unmarshal(value, &entry); err != nil {
			log.Printf("Error parsing JSON for key %s: %v", key, err)
			continue
		}

		if entry.Ver > 1 {
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Fatal(err)
	}
}
