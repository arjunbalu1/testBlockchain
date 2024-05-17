package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

// Transaction represents the structure stored in the database.
type Transaction struct {
	Key   string  `json:"key"`
	Value int     `json:"value"`
	Ver   float64 `json:"ver"`
}

func main() {
	// Open the LevelDB file.
	dbPath := filepath.Join("..", "leveldb_data") // Adjust the path to where your LevelDB data is stored.
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	// Iterate over each entry in the database.
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Use key/value.
		var txn Transaction
		if err := json.Unmarshal(iter.Value(), &txn); err != nil {
			log.Printf("Error unmarshaling data: %v", err)
			continue
		}

		// Check if version is not 1.0.
		if txn.Ver != 1.0 {
			fmt.Printf("Key: %s, Value: %d, Ver: %f\n", txn.Key, txn.Value, txn.Ver)
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}
}
