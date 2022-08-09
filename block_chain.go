/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 5:46 PM
 * @Desc:
 */

package main

import (
	"block_chain/base58"
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// BlockChain create block chain
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

func CreateBlockChain(miner string) *BlockChain {
	if IsFileExist(BlockChainDBName) {
		fmt.Printf("block chain has already exist!\n")
		return nil
	}
	db, err := bolt.Open(BlockChainDBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var tail []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		if b == nil {
			b, err = tx.CreateBucket(BucketName)
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

// NewBlockChain 返回区块链实例
func NewBlockChain() *BlockChain {
	if !IsFileExist(BlockChainDBName) {
		fmt.Printf("block chain does not exist!\n")
		return nil
	}
	db, err := bolt.Open(BlockChainDBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tail []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		if b == nil {
			fmt.Printf("buket is nil, please check!\n")
			os.Exit(1)
		}
		tail = b.Get(LastHashKey)
		return nil
	})
	return &BlockChain{
		db:   db,
		tail: tail,
	}
}

// AddBlock add block
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//矿工得到交易时，第一时间对交易进行验证
	//矿工如果不验证，即使挖矿成功，广播区块后，其它验证矿工仍然会校验每一笔交易
	validTxs := []*Transaction{}
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) {
			validTxs = append(validTxs, tx)
			fmt.Printf("交易有效： %x\n", tx.TxID)
		} else {
			fmt.Printf("发现无效交易： %x\n", tx.TxID)
			return
		}
	}

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
func (bc *BlockChain) FindMyUtoxs(publicKeyHash []byte) []UTXOInfo {
	var UTXInfos []UTXOInfo
	it := bc.NewIterator()
	//已经消耗过的
	spentUtxos := make(map[string][]int64)
	//遍历block
	for block := it.Next(); block != nil; block = it.Next() {
		//遍历交易
		for _, transaction := range block.Transactions {
			//遍历input
			if transaction.IsCoinbase() == false { //普通交易才需要遍历
				for _, input := range transaction.TxInputs {
					//找到属于我的所有output
					if bytes.Equal(HashPublicKey(input.PublicKey), publicKeyHash) {
						fmt.Printf("%s find my input, i: %d\n", publicKeyHash, input.Index)
						key := string(input.TxID)
						spentUtxos[key] = append(spentUtxos[key], input.Index)
					}
				}
			}

			key := string(transaction.TxID)
			indexes := spentUtxos[key]
		OUTPUT:
			//遍历output
			for i, output := range transaction.TxOutputs {
				//找到属于我的所有output
				if bytes.Equal(publicKeyHash, output.PublicKeyHash) {
					if len(indexes) != 0 {
						fmt.Printf("当前笔交易中又被消耗过的output\n")
						for _, j := range indexes {
							if int64(i) == j {
								fmt.Printf("i==j,当前的output已经被消耗过了，跳过不统计\n")
								continue OUTPUT
							}
						}
					}
					utxoInfo := UTXOInfo{
						TxID:   transaction.TxID,
						Index:  int64(i),
						Output: output,
					}
					UTXInfos = append(UTXInfos, utxoInfo)
					fmt.Printf("%s find my out, i: %d\n", publicKeyHash, i)
				}
			}
		}
	}
	return UTXInfos
}
func (bc *BlockChain) GetBalance(address string) float64 {
	decodeInfo := base58.Decode(address)
	publicKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	utxoInfos := bc.FindMyUtoxs(publicKeyHash)
	var total float64
	for _, utxoInfo := range utxoInfos {
		total += utxoInfo.Output.Value
	}
	fmt.Printf(" %s balance: %f\n", address, total)
	return total
}
func (bc *BlockChain) FindNeedUtxos(publicKeyHash []byte, amount float64) (map[string][]int64, float64) {
	var resValue float64
	var needUtxos = make(map[string][]int64)
	//复用FindMyUtxo函数，这个函数包含所有的信息
	//decodeInfo := base58.Decode(from)
	//publicKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	utxoInfos := bc.FindMyUtoxs(publicKeyHash)
	for _, utxoInfo := range utxoInfos {
		key := string(utxoInfo.TxID)
		needUtxos[key] = append(needUtxos[key], utxoInfo.Index)
		resValue += utxoInfo.Output.Value
		if resValue >= amount {
			break
		}
	}
	return needUtxos, resValue
}

type UTXOInfo struct {
	TxID   []byte
	Index  int64    //output索引值
	Output TxOutput //output本身
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//1.遍历账本，找到所有引用的交易
	prevTxs := make(map[string]Transaction)
	//遍历tx的input，通过id去查找所引用的交易
	for _, input := range tx.TxInputs {
		prevTx := bc.FindTransaction(input.TxID)
		if prevTx == nil {
			fmt.Printf("没有找到交易，这是一个错误： %x\n", input.TxID)
			continue
		} else {
			prevTxs[string(input.TxID)] = *prevTx
		}

	}
	tx.Sign(privateKey, prevTxs)
}
func (bc *BlockChain) FindTransaction(txID []byte) *Transaction {
	//遍历账本，通过对比ID来识别
	it := bc.NewIterator()
	for block := it.Next(); block != nil; block = it.Next() {
		for _, transaction := range block.Transactions {
			if bytes.Equal(txID, transaction.TxID) {
				fmt.Printf("找到了交易: %s\n", string(transaction.TxID))
				return transaction
			}

		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	fmt.Printf("没有找到了交易\n")
	return nil
}

// VerifyTransaction 矿工校验流程
//1。找到交易input所引用的所有的交易
//2。对交易进行签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() { //挖矿交易直接返回true，
		return true
	}
	//1.遍历账本，找到所有引用的交易
	prevTxs := make(map[string]Transaction)
	//遍历tx的input，通过id去查找所引用的交易
	for _, input := range tx.TxInputs {
		prevTx := bc.FindTransaction(input.TxID)
		if prevTx == nil {
			fmt.Printf("没有找到交易，这是一个错误： %x\n", input.TxID)
			continue
		} else {
			prevTxs[string(input.TxID)] = *prevTx
		}

	}
	return tx.Verify(prevTxs)
}
