/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/1 6:12 PM
 * @Desc:
 */

package main

import "fmt"

func (cli *CLI) AddBlock(txs []*Transaction) {
	bc := NewBlockChain()
	if bc != nil {
		defer bc.db.Close()
	} else {
		fmt.Printf("block chain is nil\n")
		return
	}
	bc.AddBlock(txs)
	fmt.Printf("add block ok\n")
}
func (cli *CLI) PrintChain() {
	bc := NewBlockChain()
	if bc != nil {
		defer bc.db.Close()
	} else {
		fmt.Printf("block chain is nil\n")
		return
	}
	blockIterator := bc.NewIterator()
	for block := blockIterator.Next(); block != nil; block = blockIterator.Next() {
		fmt.Printf("************************************\n")
		fmt.Printf("pre block hash: %x\n", block.PrevBlockHash)
		fmt.Printf("nonce: %x\n", block.Nonce)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("transactions: %+v", block.Transactions)
		pow := NewProofOfWork(block)
		fmt.Printf("is valid: %+v\n", pow.IsValid())
		fmt.Printf("====================================\n")
	}
}

func (cli *CLI) GetBalance(address string) {
	bc := NewBlockChain()
	if bc != nil {
		defer bc.db.Close()
	} else {
		fmt.Printf("block chain is nil\n")
		return
	}
	bc.GetBalance(address)
}
func (cli *CLI) Send(from, to string, amount float64, miner string) {
	bc := NewBlockChain()
	if bc != nil {
		defer bc.db.Close()
	} else {
		fmt.Printf("block chain is nil\n")
		return
	}
	//创建挖矿交易
	coinbase := NewCoinbaseTx(miner, "hello world")
	//创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	txs := []*Transaction{coinbase}
	if tx != nil {
		txs = append(txs, tx)
	}
	if tx == nil {
		fmt.Printf("无效交易")
		return
	}
	//添加区块

	bc.AddBlock(txs)
	fmt.Printf("wakuang chenggong!")
}
func (cli *CLI) CreateBlockChain(address string) {
	bc := CreateBlockChain(address)
	if bc != nil {
		defer bc.db.Close()
	} else {
		fmt.Printf("block chain is nil\n")
		return
	}
	fmt.Printf("create block chain success!\n")
}
func (cli *CLI) CreateWallet(address string) {
	w := NewWalletKeyPair()
	fmt.Printf("wallet: %s\n", w.GetAddress())
}
