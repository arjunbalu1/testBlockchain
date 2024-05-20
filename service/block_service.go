package service

import (
	"blockchain/db"
	"blockchain/model"
	"blockchain/utils"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type BlockService struct {
	db              *db.LevelDB
	blockQueue      chan *model.Block
	maxTxnsPerBlock int
	lastBlockHash   string // Store the hash of the last block
	currentBlock    *model.Block
	mu              sync.Mutex
	timer           *time.Timer
}

func NewBlockService(ldb *db.LevelDB, queue chan *model.Block, maxTxns int) *BlockService {
	// Load the last block number from file
	lastBlockNumber := GetLastBlockNumber()

	bs := &BlockService{
		db:              ldb,
		blockQueue:      queue,
		maxTxnsPerBlock: maxTxns,
		currentBlock:    model.NewBlock(lastBlockNumber+1, ""), // Start with the next block number
		timer:           time.NewTimer(15 * time.Second),
	}
	go bs.commitBlockPeriodically()
	return bs
}

func (bs *BlockService) ProcessTransactions(txns []json.RawMessage) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	for _, txn := range txns {
		if len(bs.currentBlock.Txns) >= bs.maxTxnsPerBlock {
			bs.commitBlock()
			bs.currentBlock = model.NewBlock(bs.currentBlock.BlockNumber+1, bs.lastBlockHash)
		}

		var t model.Transaction
		if err := json.Unmarshal(txn, &t); err != nil {
			fmt.Println("Failed to unmarshal transaction:", err)
			continue
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
		bs.currentBlock.Txns = append(bs.currentBlock.Txns, t)
	}

	if len(bs.currentBlock.Txns) >= bs.maxTxnsPerBlock {
		bs.commitBlock()
		bs.currentBlock = model.NewBlock(bs.currentBlock.BlockNumber+1, bs.lastBlockHash)
	}

	bs.resetTimer()
}

func (bs *BlockService) commitBlockPeriodically() {
	for range bs.timer.C {
		bs.mu.Lock()
		if len(bs.currentBlock.Txns) > 0 {
			bs.commitBlock()
			bs.currentBlock = model.NewBlock(bs.currentBlock.BlockNumber+1, bs.lastBlockHash)
		}
		bs.mu.Unlock()
		time.Sleep(1 * time.Second)
		fmt.Println("Exiting due to inactivity...")
		os.Exit(0)
	}
}

func (bs *BlockService) commitBlock() {
	bs.currentBlock.UpdateStatusToCommitted()
	blockHash, err := utils.GenerateHash(bs.currentBlock)
	if err != nil {
		fmt.Printf("Failed to generate hash for the block: %v\n", err)
		return
	}
	bs.lastBlockHash = blockHash
	fmt.Printf("Committing block %d\n", bs.currentBlock.BlockNumber)
	bs.blockQueue <- bs.currentBlock

	// Save the last block number to file
	SaveLastBlockNumber(bs.currentBlock.BlockNumber)
}

func (bs *BlockService) resetTimer() {
	if !bs.timer.Stop() {
		<-bs.timer.C
	}
	bs.timer.Reset(15 * time.Second)
	fmt.Println("Timer reset")
}
