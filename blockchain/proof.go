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

//Proof of work:Force network to view work to add a block to the chain
//Validation of this proof
//Concept: work must be hard to do but validation must be trivial

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s
//(if difficulty goes uo there must be more preceding zeroes)

//here difficulty is kept to be constant but an actual implementation of blockchain, difficulty goes up

const Difficulty = 12

//Traget- A number that is the requirement that is derieved from diffculty
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

//takes in pointer to a block and creates pointer to a pow
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	//left shift by that in target(number of bytes in a hash - difficulty)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

//adding Hash transactions
//Take the block and create a new hash using PrevHash, transaction, nonce and difficulty
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
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
		//prepare data
		data := pow.InitData(nonce)
		//hash data
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		//convert Hash into bigint
		intHash.SetBytes(hash[:])
		//compare this with por target
		//if it is -1, the we have signed the block
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:]
}

// validation function(algorithm relatively easy)
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

//Take integer and output a slice of bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
