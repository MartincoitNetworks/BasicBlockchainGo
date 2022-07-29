package blockchain

import "github.com/dgraph-io/badger"

//as our blockchain is in a db we need an iterator
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//Take blockchain struct and create a blockchain iterator struct
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

//Iterate through our blockchain from last added to genesis block
func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
