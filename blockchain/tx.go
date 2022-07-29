package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/tensor-programming/golang-blockchain/wallet"
)

//value is the value in tokens assigned in the blockchain
//pubkey is what is used to unlock the tokens in the value
//inputs and outputs are indivisible
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

// TxInput references the previous outputs
// ID referenced the transaction that the output is a part of
//index of transactions. For example if transactions has 3 outputs
//sig is similar to pubkey:users account (but is actually a complex)
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)
	return outputs
}
