package blockchain

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "/tmp/blocks" // Path where BadgerDB will store blockchain blocks
)

type Blockchain struct {
	LastHash []byte // Stores the hash of the last block added to the blockchain.
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *Blockchain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error { // executes a write transaction on the database.
		// lh → lasthash
		// txn.Get([]byte("lh")) tries to fetch the key "lh" from the database.
		// and If the key is not found (badger.ErrKeyNotFound), it means the blockchain does not exist yet.
		if _, err := txn.Get([]byte("lh")); errors.Is(err, badger.ErrKeyNotFound) {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis block")
			// genesis.Hash → Key to access the block.
			// genesis.Serialize() → Block converted to bytes.
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			// Saves the hash of the genesis block under the key "lh" to indicate that it is the last block added.
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			// txn.Get([]byte("lh")) retrieves the hash of the last stored block.
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			// item.Value() accesses the value stored at key "lh", which is the hash of the last block.
			error := item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			Handle(error)
			return err
		}
	})
	Handle(err)
	return &Blockchain{lastHash, db}
}

// AddBlock wse are adding a new blocks to an existing blockchain.
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		error := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})
		Handle(error)

		return err
	})
	Handle(err)
	newBlock := CreateBlock(data, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var blockData []byte
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		error := item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)
			return nil
		})
		Handle(error)
		block = DeserializeBlock(blockData)
		return err
	})
	Handle(err)
	iter.CurrentHash = block.PrevHash
	return block
}
