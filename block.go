package main

import (
	"bytes"
	"encoding/gob"
)

type Block struct {
	Id            string
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     int64
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	logError(err)
	return result.Bytes()
}

func newBlock(data string, prevBlockHash []byte) *Block {
	id := getUUID()
	timestamp := getTimestamp()
	b := &Block{Id: id, Data: []byte(data), Timestamp: timestamp, PrevBlockHash: prevBlockHash}
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()
	b.Hash = hash
	b.Nonce = nonce
	return b
}

func genesisBlock() *Block {
	return newBlock("Genesis Block", []byte{})
}

func Deserialize(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	logError(err)
	return &block
}
