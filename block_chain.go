/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 5:46 PM
 * @Desc:
 */

package main

import (
	"fmt"
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

func NewBlockChain(miner string) *BlockChain {
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
			//创世快中只有一个挖矿交易，只有Coinbase
			coinbase := NewCoinbaseTx(miner, "fist block")
			genesisBlock := NewBlock([]*Transaction{coinbase}, FirstBlock)
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
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		if b == nil {
			panic("bucket nil")
		}
		block := NewBlock(txs, bc.tail)
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
		if len(blockInfo) == 0 {
			return nil
		}
		block = Deserialize(blockInfo)
		it.current = block.PrevBlockHash
		return nil
	})
	if err != nil {
		panic(err)
	}
	return block
}
func (bc *BlockChain) FindMyUtoxs(address string) []TxOutput {
	//todo
	it := bc.NewIterator()
	var UtxOutputs []TxOutput
	//已经消耗过的
	spentUtxos := make(map[string][]int64)
	//遍历block
	for block := it.Next(); block != nil; block = it.Next() {
		//遍历交易
		for _, transaction := range block.Transactions {
			//遍历input
			for _, input := range transaction.TxInputs {
				//找到属于我的所有output
				if address == input.Address {
					fmt.Printf("%s find my input, i: %d\n", address, input.Index)
					key := string(input.TxID)
					spentUtxos[key] = append(spentUtxos[key], input.Index)
				}
			}
		OUTPUT:
			//遍历output
			for i, output := range transaction.TxOutputs {
				//找到属于我的所有output
				if address == output.Address {

					key := string(transaction.TxID)
					indexes := spentUtxos[key]
					if len(indexes) != 0 {
						fmt.Printf("当前笔交易中又被消耗过的output\n")
						for _, j := range indexes {
							if int64(i) == j {
								fmt.Printf("i==j,当前的output已经被消耗过了，跳过不统计\n")
								continue OUTPUT
							}
						}
					}

					UtxOutputs = append(UtxOutputs, output)
					fmt.Printf("%s find my out, i: %d\n", address, i)
				}
			}
		}
	}
	return UtxOutputs
}
func (bc *BlockChain) GetBalance(address string) float64 {
	utxos := bc.FindMyUtoxs(address)
	var total float64
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf(" %s balance: %f\n", address, total)
	return total
}
func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {
	//todo 正道utos集合
	it := bc.NewIterator()
	//var UtxOutputs []TxOutput
	var resValue float64
	var needUtxos = make(map[string][]int64)
	//已经消耗过的
	spentUtxos := make(map[string][]int64)
	//遍历block
	for block := it.Next(); block != nil; block = it.Next() {
		//遍历交易
		for _, transaction := range block.Transactions {
			//遍历input
			for _, input := range transaction.TxInputs {
				//找到属于我的所有input
				if from == input.Address {
					fmt.Printf("%s find my input, i: %d\n", from, input.Index)
					key := string(input.TxID)
					spentUtxos[key] = append(spentUtxos[key], input.Index)
				}
			}
		OUTPUT:
			//遍历output
			for i, output := range transaction.TxOutputs {
				//找到属于我的所有output
				if from == output.Address {

					key := string(transaction.TxID)
					indexes := spentUtxos[key]
					if len(indexes) != 0 {
						fmt.Printf("当前笔交易中又被消耗过的output\n")
						for _, j := range indexes {
							if int64(i) == j {
								fmt.Printf("i==j,当前的output已经被消耗过了，跳过不统计\n")
								continue OUTPUT
							}
						}
					}

					//UtxOutputs = append(UtxOutputs, output)\
					needUtxos[key] = append(needUtxos[key], int64(i))
					resValue += output.Value
					if resValue >= amount {
						return needUtxos, resValue
					}
					fmt.Printf("%s find my out, i: %d\n", from, i)
				}
			}
		}
	}
	return nil, 0
}
