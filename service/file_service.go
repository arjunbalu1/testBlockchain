package service

import (
	"blockchain/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type FileService struct {
	filePath string
}

func NewFileService(filePath string) *FileService {
	return &FileService{filePath: filePath}
}

// WriteBlockToFile writes a single block to a JSON file, maintaining it as part of an array of blocks
func (fs *FileService) WriteBlockToFile(block *model.Block) error {
	// Open the file in read-write mode, create it if it does not exist
	file, err := os.OpenFile(fs.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read the existing contents into a buffer
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Determine whether the file already contains a JSON array
	var blocks []*model.Block
	if err = json.Unmarshal(data, &blocks); err != nil && len(data) != 0 {
		return fmt.Errorf("error parsing existing JSON: %v", err)
	}

	// Append the new block to the array of blocks
	blocks = append(blocks, block)

	// Marshal the updated blocks array into JSON
	updatedData, err := json.Marshal(blocks)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Truncate the file and write the updated JSON
	if err = file.Truncate(0); err != nil {
		return fmt.Errorf("error truncating file: %v", err)
	}
	if _, err = file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking file: %v", err)
	}
	if _, err = file.Write(updatedData); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	fmt.Printf("Block %d written to file\n", block.BlockNumber)

	return nil
}

// FetchBlockByNumber fetches a single block by its number from the file
func (fs *FileService) FetchBlockByNumber(blockNumber uint64) (*model.Block, error) {
	file, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, err
	}
	var blocks []*model.Block
	if err := json.Unmarshal(file, &blocks); err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.BlockNumber == blockNumber {
			return block, nil
		}
	}
	return nil, fmt.Errorf("block number %d not found", blockNumber)
}

// FetchAllBlocks retrieves all blocks from the file
func (fs *FileService) FetchAllBlocks() ([]*model.Block, error) {
	file, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, err
	}
	var blocks []*model.Block
	if err := json.Unmarshal(file, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

func GetLastBlockNumber() uint64 {
	data, err := ioutil.ReadFile("last_block_number.txt")
	if err != nil {
		return 0 // Start from block number 0 because increment happens before processing
	}
	lastBlockNumber, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0
	}
	return lastBlockNumber
}

func SaveLastBlockNumber(blockNumber uint64) {
	ioutil.WriteFile("last_block_number.txt", []byte(fmt.Sprintf("%d", blockNumber)), 0644)
	fmt.Printf("Last block number %d saved\n", blockNumber)
}
