/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/2 4:18 PM
 * @Desc:
 */

package main

import (
	"block_chain/base58"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//交易结构
type TxInput struct {
	TxID  []byte //transaction id
	Index int64  //output
	//Address string //解锁脚本，先使用地址模拟
	Signature []byte //交易签名
	PublicKey []byte //公钥本身
}
type TxOutput struct {
	Value float64 //money
	//Address string  //锁定脚本
	PublicKeyHash []byte //公钥的hash，不是公钥
}

// Lock 给定一个地址，得到这个地址的公钥hash
func (output *TxOutput) Lock(address string) {
	//address-》public key hash
	decodeInfo := base58.Decode(address)
	publicKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	output.PublicKeyHash = publicKeyHash
}
func NewTxOutput(value float64, address string) *TxOutput {
	output := &TxOutput{
		Value:         value,
		PublicKeyHash: nil,
	}
	output.Lock(address)
	return output
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

const Reward = 12.5

// NewCoinbaseTx 挖矿
//特点：只有输出，没有输入
func NewCoinbaseTx(miner string, data string) *Transaction {
	//todo
	var inputs []TxInput

	inputs = append(inputs, TxInput{
		TxID:      nil,
		Index:     -1,
		Signature: nil,
		PublicKey: []byte(data),
	})

	var outputs []TxOutput
	outputs = append(outputs, *NewTxOutput(Reward, miner))
	transaction := &Transaction{
		TxID:      nil,
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	transaction.SetTxID()
	return transaction
}

func (tx *Transaction) IsCoinbase() bool {
	//特点：只有一个input；引用ID是nil；引用的索引是-1
	if len(tx.TxInputs) == 1 && tx.TxInputs[0].TxID == nil && tx.TxInputs[0].Index == -1 {
		return true
	}
	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1打开钱包
	ws := NewWallets()

	//获取密钥对
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		return nil
	}
	//2。获取公钥私钥
	//pivateKey := wallet.PrivateKey, 目前用不到
	publicKey := wallet.PublicKey

	publicKeyHash := HashPublicKey(wallet.PublicKey)

	utxos := make(map[string][]int64)
	var resValue float64
	//假如李四转赵六4元钱，返回的信息为：
	//utxos[0x333]=int64{0,1}
	utxos, resValue = bc.FindNeedUtxos(publicKeyHash, amount)

	if resValue < amount {
		fmt.Printf("less money \n")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput
	for txid, indexes := range utxos {
		for _, i := range indexes {
			input := TxInput{
				TxID:      []byte(txid),
				Index:     i,
				Signature: nil,
				PublicKey: publicKey,
			}
			inputs = append(inputs, input)
		}
	}

	output := NewTxOutput(amount, to)

	outputs = append(outputs, *output)
	if resValue > amount {
		output1 := NewTxOutput(resValue-amount, from)
		outputs = append(outputs, *output1)
	}
	tx := Transaction{
		TxID:      nil,
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	tx.SetTxID()
	return &tx
}
