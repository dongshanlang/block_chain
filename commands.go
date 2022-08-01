/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/1 6:12 PM
 * @Desc:
 */

package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Printf("add block ok\n")
}
func (cli *CLI) PrintChain() {
	blockIterator := cli.bc.NewIterator()
	for block := blockIterator.Next(); block != nil; block = blockIterator.Next() {
		fmt.Printf("pre block hash: %x\n", block.PrevBlockHash)
		fmt.Printf("nonce: %x\n", block.Nonce)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Println("data: ", string(block.Data))
		pow := NewProofOfWork(block)
		fmt.Printf("is valid: %+v\n", pow.IsValid())
	}
}
func printBlockChain(bc *BlockChain) {
	blockIterator := bc.NewIterator()
	for block := blockIterator.Next(); block != nil; block = blockIterator.Next() {
		fmt.Printf("pre block hash: %x\n", block.PrevBlockHash)
		fmt.Printf("nonce: %x\n", block.Nonce)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Println("data: ", string(block.Data))
		pow := NewProofOfWork(block)
		fmt.Printf("is valid: %+v\n", pow.IsValid())
	}
}
