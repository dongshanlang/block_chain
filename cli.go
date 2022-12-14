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
	"strconv"
)

const (
	Usage = `
	./blockchain printChain                           打印区块链
	./blockchain getBalance   address                 获取余额
	./blockchain send FROM TO AMOUNT MINER            转账命令
	./blockchain createBlockChain address             创建区块链
	./blockchain createWallet                         创建钱包
	./blockchain listAddress                          打印地址
	./blockchain printTX                              打印所有交易
`
)

type CLI struct {
	//bc *BlockChain，自己获取区块链实例
}

func (cli *CLI) Run() {
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Printf("%s\n", Usage)
		return
	}
	switch cmds[1] {
	case "printChain":
		cli.PrintChain()
	case "getBalance":
		cli.GetBalance(cmds[2])
	case "createBlockChain":
		if len(cmds) != 3 {
			fmt.Printf("无效参数\n")
			fmt.Printf("%s\n", Usage)
			return
		}
		address := cmds[2]
		cli.CreateBlockChain(address)
	case "send":
		fmt.Printf("send\n")
		if len(cmds) != 6 {
			fmt.Printf("无效参数\n")
			fmt.Printf("%s\n", Usage)
			return
		}
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		cli.Send(from, to, amount, miner)
	case "createWallet":
		if len(cmds) != 2 {
			fmt.Printf("无效参数\n")
			fmt.Printf("%s\n", Usage)
			return
		}
		//address := cmds[2]
		cli.CreateWallet("address")
	case "listAddress":
		fmt.Printf("listAddress\n")
		if len(cmds) != 2 {
			fmt.Printf("无效参数\n")
			fmt.Printf("%s\n", Usage)
			return
		}
		cli.ListAddress()
	case "printTX":
		cli.PrintTX()

	default:
		fmt.Printf("%s\n", Usage)
	}
}
