package service

import (
	"blockchain/db"
	"blockchain/model"
	"blockchain/utils" // Ensure utils is imported to use GenerateHash
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type BlockService struct {
	db              *db.LevelDB
	blockQueue      chan *model.Block
	maxTxnsPerBlock int
	lastBlockHash   string // Store the hash of the last block
}

func NewBlockService(ldb *db.LevelDB, queue chan *model.Block, maxTxns int) *BlockService {
	return &BlockService{
		db:              ldb,
		blockQueue:      queue,
		maxTxnsPerBlock: maxTxns,
	}
}

func (bs *BlockService) ProcessTransactions(txns []json.RawMessage, blockNum uint64) {
	block := model.NewBlock()
	block.BlockNumber = blockNum
	block.Timestamp = time.Now().Unix()
	block.PreviousHash = bs.lastBlockHash // Set the previous hash
	block.BlockStatus = model.Pending

	var wg sync.WaitGroup
	txnsToProcess := make([]model.Transaction, 0, bs.maxTxnsPerBlock)
	for i, txn := range txns {
		if i >= bs.maxTxnsPerBlock {
			break
		}
		wg.Add(1)
		go func(txn json.RawMessage) {
			defer wg.Done()
			var t model.Transaction
			if err := json.Unmarshal(txn, &t); err != nil {
				fmt.Println("Failed to unmarshal transaction:", err)
				return
			}

			data, err := bs.db.Get(t.Key)
			if err != nil {
				fmt.Println("Failed to get data from LevelDB for key:", t.Key)
				t.Valid = false
			} else {
				var currentVal model.Transaction
				json.Unmarshal(data, &currentVal)
				if t.Ver == currentVal.Ver {
					t.Valid = true
					t.Ver += 1.0 // Increment version
					updatedVal := model.Transaction{
						Key:   t.Key,
						Value: t.Value,
						Ver:   t.Ver,
						Valid: true,
					}
					valBytes, err := json.Marshal(updatedVal)
					if err == nil {
						err = bs.db.Put(t.Key, valBytes)
						if err != nil {
							fmt.Printf("Error updating key %s in LevelDB: %v\n", t.Key, err)
						}
					} else {
						fmt.Printf("Error marshaling updated transaction for key %s: %v\n", t.Key, err)
					}
				} else {
					t.Valid = false
				}
			}

			t.Hash, _ = utils.GenerateHash(t)
			txnsToProcess = append(txnsToProcess, t)
		}(txn)
	}
	wg.Wait()

	block.PushValidTxns(txnsToProcess)
	block.UpdateStatusToCommitted()
	blockHash, err := utils.GenerateHash(block) // Directly use GenerateHash
	if err != nil {
		fmt.Printf("Failed to generate hash for the block: %v\n", err)
		return
	}
	bs.lastBlockHash = blockHash // Update the last block hash
	bs.blockQueue <- block
}
