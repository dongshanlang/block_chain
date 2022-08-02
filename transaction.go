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
func NewCoinbaseTx(miner string) *Transaction {
	//todo
	var inputs []TxInput

	inputs = append(inputs, TxInput{
		TxID:    nil,
		Index:   -1,
		Address: gensisInfo,
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
