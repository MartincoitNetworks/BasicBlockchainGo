package blockchain

//Not everyone wanted to have a full node on thier local system
//SPVs are a way around that. Nodes that do not verify or download entire blockchain. It finds transactions in bloack and is linked to a full node
//Merkle trees are used to obtain transaction hashes and are saved inside a blockchain header
//leaves of merkle tree are even
//we can just look at the root of the tree to see if a certain transaction is available in a block

import (
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

//recursive tree structure
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

//create new merkle node

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}
	//check if sides of merkle tree exists
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}

//recursively add the hashs and create new hashes from bottom to top of tree
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	if len(nodes) == 0 {
		log.Panic("No merkel nodes")
	}

	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		var level []MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			node := NewMerkleNode(&nodes[i], &nodes[i+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	tree := MerkleTree{&nodes[0]}

	return &tree
}
