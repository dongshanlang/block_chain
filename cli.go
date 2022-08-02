/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/1 6:02 PM
 * @Desc:
 */

package main

import (
	"fmt"
	"os"
)

const (
	Usage = `
	./blockchain addBlock "xxxxxxxx" 添加数据到区块链
	./blockchain printChain         打印区块链
`
)

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) Run() {
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Printf("%s\n", Usage)
		return
	}
	switch cmds[1] {
	case "addBlock":
		if len(cmds) != 3 {
			fmt.Printf("%s\n", Usage)
			return
		}
		//cli.AddBlock(cmds[2])//todo
	case "printChain":
		cli.PrintChain()
	default:
		fmt.Printf("%s\n", Usage)
	}
}
