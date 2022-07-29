package blockchain

//Blockchain- a public database that is distributed amongst multiple different peers.
//what makes it unique is it does not rely on trusting nodes
import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//Block-contains data in a database
//Hash- Hash of current block
//Data- Anything from images to ledgers stored in a block
//PrevHash - Last blocks hash(used for connecting blocks like a linkedlist)
//each block must have a minimum of one transaction
type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

//function to include transactions in our proof of work
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	//representing our blocks using merkle tree
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

//Creating a block
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {

	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	// run proof of work algo on each block
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//The first block in the chain as it does not have a prevHash
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

//Badger DB only allows arrays of bytes we need to serialise and deserialise block data structure
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

//Handle function for errors
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
