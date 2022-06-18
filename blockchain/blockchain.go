package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

const (
	lhKey = "lh"
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block
	err := i.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.CurrentHash)
		err = item.Value(func(encodedBlock []byte) error {
			block = Deserialize(encodedBlock)
			return nil
		})
		return err
	})
	if err != nil {
		log.Panicln(err)
	}
	i.CurrentHash = block.PrevHash
	return block
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts);
	if err != nil {
		log.Panicln(err)
	}
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(lhKey)); err == badger.ErrKeyNotFound {
			fmt.Println("No genesis block found")
			gen := Genesis()
			fmt.Println("Init genesis block")
			err = txn.Set(gen.Hash, gen.Serialize())
			err = txn.Set([]byte(lhKey), gen.Hash)
			lastHash = gen.Hash
		} else {
			item, err := txn.Get([]byte(lhKey))
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
		return err
	})
	if err != nil {
		log.Panicln(err)
	}

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(lhKey))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	if err != nil {
		log.Panicln(err)
	}

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte(lhKey), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	if err != nil {
		log.Panicln(err)
	}
}
