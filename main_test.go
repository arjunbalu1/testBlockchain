package main

import (
	"blockchain/db"
	"blockchain/model"
	"blockchain/service"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupLevelDB(t *testing.T) *db.LevelDB {
	path := filepath.Join(".", "test_leveldb_data")
	ldb := db.NewLevelDB(path)
	ldb.InitializeDB()
	return ldb
}

func cleanupLevelDB(ldb *db.LevelDB) {
	ldb.Close()
	os.RemoveAll("test_leveldb_data")
}

func TestInitializeDB(t *testing.T) {
	ldb := setupLevelDB(t)
	defer cleanupLevelDB(ldb)

	for i := 1; i <= 1000; i++ {
		key := fmt.Sprintf("SIM%d", i)
		data, err := ldb.Get(key)
		assert.NoError(t, err)
		var entry map[string]interface{}
		err = json.Unmarshal(data, &entry)
		assert.NoError(t, err)
		assert.Equal(t, float64(i), entry["val"].(float64))
		assert.Equal(t, 1.0, entry["ver"].(float64))
	}
}

func TestProcessTransactions(t *testing.T) {
	ldb := setupLevelDB(t)
	defer cleanupLevelDB(ldb)

	blockQueue := make(chan *model.Block, 100)
	blockService := service.NewBlockService(ldb, blockQueue, 10)

	transactions := []json.RawMessage{
		json.RawMessage(`{"key": "SIM1", "value": 2, "ver": 1.0}`),
		json.RawMessage(`{"key": "SIM2", "value": 3, "ver": 1.0}`),
		json.RawMessage(`{"key": "SIM3", "value": 4, "ver": 2.0}`),
	}

	blockService.ProcessTransactions(transactions)

	time.Sleep(1 * time.Second)

	block := <-blockQueue

	assert.Equal(t, uint64(0), block.BlockNumber)
	assert.Equal(t, 2, block.Txns[0].Value)
	assert.Equal(t, 2.0, block.Txns[0].Ver)
	assert.True(t, block.Txns[0].Valid)

	assert.Equal(t, 3, block.Txns[1].Value)
	assert.Equal(t, 2.0, block.Txns[1].Ver)
	assert.True(t, block.Txns[1].Valid)

	assert.Equal(t, 4, block.Txns[2].Value)
	assert.Equal(t, 2.0, block.Txns[2].Ver)
	assert.False(t, block.Txns[2].Valid)
}

func TestBlockCommitment(t *testing.T) {
	ldb := setupLevelDB(t)
	defer cleanupLevelDB(ldb)

	blockQueue := make(chan *model.Block, 100)
	blockService := service.NewBlockService(ldb, blockQueue, 2)

	transactions := []json.RawMessage{
		json.RawMessage(`{"key": "SIM1", "value": 2, "ver": 1.0}`),
		json.RawMessage(`{"key": "SIM2", "value": 3, "ver": 1.0}`),
	}

	blockService.ProcessTransactions(transactions)
	time.Sleep(1 * time.Second)

	block := <-blockQueue
	assert.Equal(t, uint64(0), block.BlockNumber)

	extraTransactions := []json.RawMessage{
		json.RawMessage(`{"key": "SIM4", "value": 5, "ver": 1.0}`),
	}
	blockService.ProcessTransactions(extraTransactions)
	time.Sleep(1 * time.Second)

	block = <-blockQueue
	assert.Equal(t, uint64(1), block.BlockNumber)
}

func TestFileOperations(t *testing.T) {
	filePath := filepath.Join(".", "test_blocks.json")
	fileService := service.NewFileService(filePath)

	block := &model.Block{
		BlockNumber:  1,
		Txns:         []model.Transaction{},
		Timestamp:    time.Now().Unix(),
		BlockStatus:  model.Committed,
		PreviousHash: "0xabc123",
	}

	err := fileService.WriteBlockToFile(block)
	assert.NoError(t, err)

	fetchedBlock, err := fileService.FetchBlockByNumber(1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), fetchedBlock.BlockNumber)

	allBlocks, err := fileService.FetchAllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allBlocks))

	os.Remove(filePath)
}

func TestBlockProcessingTime(t *testing.T) {
	ldb := setupLevelDB(t)
	defer cleanupLevelDB(ldb)

	blockQueue := make(chan *model.Block, 100)
	blockService := service.NewBlockService(ldb, blockQueue, 10)

	startTime := time.Now()
	transactions := []json.RawMessage{
		json.RawMessage(`{"key": "SIM1", "value": 2, "ver": 1.0}`),
	}

	blockService.ProcessTransactions(transactions)
	time.Sleep(1 * time.Second)

	block := <-blockQueue
	duration := time.Since(startTime)

	assert.Less(t, duration.Seconds(), 15.1)      // Allow a small buffer
	assert.NotNil(t, block)                       // Ensure the block is not nil
	assert.Equal(t, uint64(0), block.BlockNumber) // Validate block number
}
