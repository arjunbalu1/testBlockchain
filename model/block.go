package model

import (
	"sync"
	"time"
)

type BlockStatus string

const (
	Committed BlockStatus = "committed"
	Pending   BlockStatus = "pending"
)

type Transaction struct {
	Key   string  `json:"key"`
	Value int     `json:"value"`
	Ver   float64 `json:"ver"`
	Valid bool    `json:"valid"`
	Hash  string  `json:"hash"`
}

type Block struct {
	BlockNumber  uint64
	Txns         []Transaction
	Timestamp    int64
	BlockStatus  BlockStatus
	PreviousHash string
	mu           sync.Mutex
}

func (b *Block) PushValidTxns(txns []Transaction) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, txn := range txns {
		b.Txns = append(b.Txns, txn)
	}
}

func (b *Block) UpdateStatusToCommitted() {
	b.mu.Lock()
	b.BlockStatus = Committed
	b.mu.Unlock()
}

func NewBlock(blockNumber uint64, previousHash string) *Block {
	return &Block{
		BlockNumber:  blockNumber,
		Txns:         []Transaction{},
		Timestamp:    time.Now().Unix(),
		BlockStatus:  Pending,
		PreviousHash: previousHash,
		mu:           sync.Mutex{},
	}
}
