package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the Data from the blocks
// Create a counter (nonce) which starts at 0
// Create a hash of the data plus the counter
// Check the hash to see if it meets a set of requirements
// :
// The first few bytes must contain 0s
// Example: 0000X123456

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,       // Hash of the previous blocks
			pow.Block.Data,           // Data stored in the blocks (ex. transactions)
			ToHex(int64(nonce)),      // Nonce converted to bytes (counter that changes with each mining attempt)
			ToHex(int64(Difficulty)), // Mining difficulty converted to bytes
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

// ToHex Receives an integer value and converts it to its binary representation
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer) // bytes.Buffer is a structure that allows you to store and manipulate binary data in memory.( It acts as a temporary container to hold the bytes of the number).
	// binary.Write writes the binary bytes of the number into the buffer
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
