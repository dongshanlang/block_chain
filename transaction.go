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
	"crypto/ecdsa"
	"crypto/rand"
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
	privateKey := wallet.PrivateKey //, 目前用不到
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

	bc.SignTransaction(&tx, privateKey)
	return &tx
}

//para1: 私钥
//para2： 交易input索引用的所有交易
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, preTXs map[string]Transaction) bool {
	fmt.Printf("对交易签名\n")
	//	1。拷贝一份交易txCopy
	//	做相应的裁剪，把每一个input的sig和publicKey设置为nil
	//	output不做改变
	txCopy := tx.TrimmedCopy()
	//	2。遍历txCopy。inputs， 把这个input所引用的output的公钥哈希拿过来，赋值给publicKey
	for i, input := range txCopy.TxInputs {
		//找到引用的交易
		preTx := preTXs[string(input.TxID)]
		output := preTx.TxOutputs[input.Index]
		txCopy.TxInputs[i].PublicKey = output.PublicKeyHash
		//input.PublicKey = output.PublicKeyHash
		//签名要对数据的hash进行签名
		//我们的数据都在交易中，我们要求交易的hash
		//Transaction的SetTxID函数就是对交易的hash
		//所以我们可以使用交易ID作为我们的签名的内容
		//	3。生成要签名的数据（哈希）
		txCopy.SetTxID()
		signData := txCopy.TxID
		txCopy.TxInputs[i].PublicKey = nil
		fmt.Printf("要签名的数据:%s\n", signData)
		//	4。对数据进行签名，r、s
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signData)
		if err != nil {
			fmt.Printf("交易签名失败：%v\n", err)
			return false
		}

		//	5。拼接r  s为字节流，赋值给原始的交易的signature字段
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TxInputs[i].Signature = signature
	}

	return true
}
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, input := range tx.TxInputs {
		inputTmp := TxInput{
			TxID:      input.TxID,
			Index:     input.Index,
			Signature: nil,
			PublicKey: nil,
		}
		inputs = append(inputs, inputTmp)

	}
	outputs = tx.TxOutputs

	transaction := Transaction{
		TxID:      tx.TxID,
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	return transaction

}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	return true
}
