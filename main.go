package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
	"time"

	"math/big"

	"github.com/google/uuid"
)

type Block struct{
	Id string
	Data []byte
	Hash []byte
	PrevBlockHash []byte
	Timestamp int64
	Nonce         int
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
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()
	b.Hash = hash
	b.Nonce = nonce
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

	
const targetBits = 24
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{block: b, target: target}
	return pow
}


func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.Data,
		pow.block.PrevBlockHash,
		[]byte(pow.block.Id),
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte(""))

	return data
}

// hardwork
func (pow *ProofOfWork) run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		// fmt.Printf("\n%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n")

	return nonce, hash[:]
}

func IntToHex(data int64) []byte {
	hex_data := fmt.Sprintf("%x", data)
	return []byte(hex_data)
}

func (pow *ProofOfWork) validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	
	return isValid
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
		pow := NewProofOfWork(block)
		isValid := pow.validate()
		fmt.Printf("PoW: %s\n", strconv.FormatBool(isValid))
		fmt.Println()
	}
}