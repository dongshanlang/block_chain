/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/2 4:18 PM
 * @Desc:
 */

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//交易结构
type TxInput struct {
	TxID    []byte //transaction id
	Index   int64  //output
	Address string //解锁脚本，先使用地址模拟
}
type TxOutput struct {
	Value   float64 //money
	Address string  //锁定脚本
}
type Transaction struct {
	TxID      []byte //transaction id
	TxInputs  []TxInput
	TxOutputs []TxOutput
}

func (tx *Transaction) SetTxID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TxID = hash[:]
}

// NewCoinbaseTx 挖矿
//特点：只有输出，没有输入
func NewCoinbaseTx(miner string, data string) *Transaction {
	//todo
	var inputs []TxInput

	inputs = append(inputs, TxInput{
		TxID:    nil,
		Index:   -1,
		Address: data,
	})

	var outputs []TxOutput
	outputs = append(outputs, TxOutput{
		Value:   12.5,
		Address: miner,
	})
	transaction := &Transaction{
		TxID:      nil,
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	transaction.SetTxID()
	return transaction
}
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	utxos := make(map[string][]int64)
	var resValue float64
	//假如李四转赵六4元钱，返回的信息为：
	//utxos[0x333]=int64{0,1}
	utxos, resValue = bc.FindNeedUtxos(from, amount)

	if resValue < amount {
		fmt.Printf("less money \n")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput
	for txid, indexes := range utxos {
		for _, i := range indexes {
			input := TxInput{
				TxID:    []byte(txid),
				Index:   i,
				Address: from,
			}
			inputs = append(inputs, input)
		}
	}

	output := TxOutput{
		Value:   amount,
		Address: to,
	}
	outputs = append(outputs, output)
	if resValue > amount {
		output1 := TxOutput{
			Value:   resValue - amount,
			Address: from,
		}
		outputs = append(outputs, output1)
	}
	tx := Transaction{
		TxID:      nil,
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	tx.SetTxID()
	return &tx
}
