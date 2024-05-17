package main

import (
	"blockchain/db"
	"blockchain/model"
	"blockchain/service"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var blockService *service.BlockService
var fileService *service.FileService
var blockQueue chan *model.Block

func init() {
	blockQueue = make(chan *model.Block, 100) // Buffered channel to prevent blocking
}

func main() {
	ldb := db.NewLevelDB(filepath.Join(".", "leveldb_data"))
	defer ldb.Close()

	ldb.InitializeDB()

	const maxTxnsPerBlock = 10 // Set the maximum number of transactions per block
	blockService = service.NewBlockService(ldb, blockQueue, maxTxnsPerBlock)
	fileService = service.NewFileService(filepath.Join(".", "blocks.json"))

	apiService := service.NewAPIService(blockService)
	router := gin.Default()
	router.POST("/transactions", apiService.PostTransactions)

	// A separate goroutine to handle blocks received from the service
	go func() {
		for block := range blockQueue {
			startTime := time.Now() // Start timing just before processing the block

			if err := fileService.WriteBlockToFile(block); err != nil {
				log.Printf("Error writing block to file: %v", err)
				continue
			}

			// Display the block processing time
			fmt.Printf("Block %d processed in %v\n", block.BlockNumber, time.Since(startTime))
		}
	}()

	log.Fatal(router.Run(":8080"))
}
