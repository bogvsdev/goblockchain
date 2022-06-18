package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, blockHash := pow.Run()
	block.Nonce = nonce
	block.Hash = blockHash[:]

	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	if err := encoder.Encode(b); err != nil {
		log.Panicln(err)
	}

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&block); err != nil {
		log.Panicln(err)
	}

	return &block
}