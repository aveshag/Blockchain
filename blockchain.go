package main

import "github.com/boltdb/bolt"

const dbFile = "blockchain.db"
const blockBucket = "block_bucket"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) AddBlocks(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	logError(err)

	newBlock := newBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		newBlockHash := newBlock.Hash
		err := b.Put(newBlockHash, newBlock.Serialize())
		logError(err)
		err = b.Put([]byte("l"), newBlockHash)
		logError(err)

		bc.tip = newBlockHash
		return nil
	})

	logError(err)

}

func NewBlockChain() *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	logError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			genesis := genesisBlock()
			b, err := tx.CreateBucket([]byte(blockBucket))
			logError(err)
			err = b.Put(genesis.Hash, genesis.Serialize())
			logError(err)
			err = b.Put([]byte("l"), genesis.Hash)
			logError(err)
			tip = genesis.Hash

		} else {
			tip = b.Get([]byte("l"))
		}

		// return nil to commit transaction
		return nil
	})
	logError(err)

	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{bc.tip, bc.db}

	return &bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	logError(err)

	i.currentHash = block.PrevBlockHash

	return block
}
