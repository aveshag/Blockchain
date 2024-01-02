package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Block struct{
	Id string
	Data []byte
	Hash []byte
	PrevBlockHash []byte
	Timestamp int64
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.Itoa(int(b.Timestamp)))
	id := []byte(b.Id)
	header := bytes.Join([][]byte{id, b.Data, b.PrevBlockHash, timestamp}, []byte(""))
	hash := sha256.Sum256(header)
	// Create array to slice
	b.Hash = hash[:]
}

func getUUID() string {
	id := uuid.New()
	return id.String()
}

func getTimestamp() int64{
	t := time.Now()
	return t.Unix()
}

func newBlock(data string, prevBlockHash []byte) *Block{
	id := getUUID()
	timestamp := getTimestamp()
	b := &Block{Id: id, Data: []byte(data), Timestamp: timestamp, PrevBlockHash: prevBlockHash}
	b.setHash()
	return b
}

type Blockchain struct{
	blocks []*Block
}


func (bc *Blockchain) addBlocks(data string) {
	prevBlock := bc.blocks[len(bc.blocks) - 1]
	newBlock := newBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func genesisBlock() *Block{
	return newBlock("Genesis Block", []byte{})
}

func getBlockchain() *Blockchain{
	return &Blockchain{blocks: []*Block{genesisBlock()}}
}

func main(){
	bc := getBlockchain()
	bc.addBlocks("Sending money to Sachin!")
	bc.addBlocks("Sending money to Rohit!")

	for _, block := range bc.blocks{
		fmt.Println("Id: ", block.Id)
		fmt.Println("Data: ", string(block.Data))
		fmt.Println("Timestamp: ", block.Timestamp)
		fmt.Printf("Hash: %X\n", block.Hash)
		fmt.Printf("Previous Block Hash: %X\n", block.PrevBlockHash)
		fmt.Println()
	}
}