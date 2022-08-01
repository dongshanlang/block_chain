/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 5:46 PM
 * @Desc:
 */

package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//create block chain
type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

var (
	BlockChainDBName = "blockChain.db"
	BucketName       = []byte("blockBucket")
	LastHashKey      = []byte("lastHashKey")
	FirstBlock       = []byte("0x0000000000000000")
)

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(BlockChainDBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tail []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(BucketName))
			if err != nil {
				log.Panic(err)
			}
			//创建创世块
			genesisBlock := NewBlock(gensisInfo, FirstBlock)
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				panic(err)
			}
			err = b.Put(LastHashKey, genesisBlock.Hash)
			if err != nil {
				panic(err)
			}
			tail = genesisBlock.Hash
		} else {
			tail = b.Get(LastHashKey)
		}
		return nil
	})
	return &BlockChain{
		db:   db,
		tail: tail,
	}
}

//add block
func (bc *BlockChain) AddBlock(data string) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		if b == nil {
			panic("bucket nil")
		}
		block := NewBlock(data, bc.tail)
		b.Put(block.Hash, block.Serialize())
		b.Put(LastHashKey, block.Hash)
		bc.tail = block.Hash
		return nil
	})
	if err != nil {
		panic(err)
	}
}

type BlockChainIterator struct {
	db      *bolt.DB
	current []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:      bc.db,
		current: bc.tail,
	}
}
func (it *BlockChainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		if b == nil {
			panic("bucket nil")
		}
		blockInfo := b.Get(it.current)
		block = Deserialize(blockInfo)
		it.current = block.PrevBlockHash
		return nil
	})
	if err != nil {
		panic(err)
	}
	return block
}
