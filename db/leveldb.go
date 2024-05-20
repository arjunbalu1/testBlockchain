package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	DB *leveldb.DB
}

func NewLevelDB(path string) *LevelDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatal("Failed to open LevelDB:", err)
	}
	return &LevelDB{DB: db}
}

func (ldb *LevelDB) Get(key string) ([]byte, error) {
	return ldb.DB.Get([]byte(key), nil)
}

func (ldb *LevelDB) Put(key string, value []byte) error {
	return ldb.DB.Put([]byte(key), value, nil)
}

func (ldb *LevelDB) Close() {
	ldb.DB.Close()
}

// InitializeDB sets up initial data in the database
func (ldb *LevelDB) InitializeDB() {
	for i := 1; i <= 1000; i++ {
		key := fmt.Sprintf("SIM%d", i)
		if data, err := ldb.Get(key); err != nil || data == nil {
			initialVal := map[string]interface{}{
				"val": i,   // Setting the value to the index number
				"ver": 1.0, // Starting version at 1.0
			}
			valBytes, err := json.Marshal(initialVal)
			if err != nil {
				log.Printf("Error marshaling initial value for %s: %v\n", key, err)
				continue
			}
			if err := ldb.Put(key, valBytes); err != nil {
				log.Printf("Error setting initial value for %s: %v\n", key, err)
			}
		}
	}
}
