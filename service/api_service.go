package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIService struct {
	BlockService *BlockService
}

func NewAPIService(blockService *BlockService) *APIService {
	return &APIService{
		BlockService: blockService,
	}
}

func (api *APIService) PostTransactions(c *gin.Context) {
	var txns []json.RawMessage
	if err := c.BindJSON(&txns); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	lastBlockNumber := GetLastBlockNumber()
	lastBlockNumber++
	SaveLastBlockNumber(lastBlockNumber)

	go func() {
		api.BlockService.ProcessTransactions(txns, lastBlockNumber)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Transactions processed successfully"})
}
