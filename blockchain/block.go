package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

type Block struct {
	Hash     []byte // Stores the hash of the current blocks.
	Data     []byte // Contains the data stored in the blocks.
	PrevHash []byte // Stores the hash of the previous blocks, connecting blocks in the chain.
	Nonce    int
}

// The CreateBlock function automates the creation of blocks on the blockchain.
// data string: Represents the data that will be stored in the blocks.
// prevHash []byte: The hash of the previous blocks, connecting the blocks in the blockchain.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// Genesis Creates and returns the genesis blocks, which is the first blocks in the blockchain.
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// Serialize Transforms a blocks into a sequence of bytes
func (b *Block) Serialize() []byte {
	var res bytes.Buffer            // where the data will be stored (data structure that stores a sequence of bytes in memory)
	encoder := gob.NewEncoder(&res) // write data to res (buffer)

	err := encoder.Encode(b) // convert the blocks into binary data and store it in res (buffer)
	Handle(err)
	return res.Bytes() // returns the bytes from the buffer, which represent a serialized blocks
}

// DeserializeBlock Convert the bytes back to a blocks
func DeserializeBlock(data []byte) *Block {
	var block Block                                  // receives the data to be deserialized
	decoder := gob.NewDecoder(bytes.NewReader(data)) // creates a reader that allows reading the bytes of the data variable

	err := decoder.Decode(&block) // reads the bytes and turns them into a blocks object
	Handle(err)
	return &block // return the reconstructed blocks
}

// Handle Simple error handling function
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
