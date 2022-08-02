/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/1 6:12 PM
 * @Desc:
 */

package main

import "fmt"

func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
	fmt.Printf("add block ok\n")
}
func (cli *CLI) PrintChain() {
	blockIterator := cli.bc.NewIterator()
	for block := blockIterator.Next(); block != nil; block = blockIterator.Next() {
		fmt.Printf("************************************\n")
		fmt.Printf("pre block hash: %x\n", block.PrevBlockHash)
		fmt.Printf("nonce: %x\n", block.Nonce)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Println("transactions: ", block.Transactions)
		pow := NewProofOfWork(block)
		fmt.Printf("is valid: %+v\n", pow.IsValid())
		fmt.Printf("====================================\n")
	}
}
