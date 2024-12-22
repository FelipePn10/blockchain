package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte // Stores the hash of the current block.
	Data     []byte // Contains the data stored in the block.
	PrevHash []byte // Stores the hash of the previous block, connecting blocks in the chain.
}

// Generate the block hash based on the data (Data)
// and the hash of the previous block (PrevHash).
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// The CreateBlock function automates the creation of blocks on the blockchain.
// data string: Represents the data that will be stored in the block.
// prevHash []byte: The hash of the previous block, connecting the blocks in the blockchain.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// We are adding a new block to an existing blockchain.
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

// Creates and returns the genesis block, which is the first block in the blockchain.
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}
