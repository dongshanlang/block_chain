/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/1 10:43 AM
 * @Desc:
 */

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const bits = 16

type ProofOfWork struct {
	block *Block

	//存储hash，借用big.Int内置的方法，Cmp
	//SetBytes
	//SetString
	target *big.Int //系统提供
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block:  block,
		target: nil,
	}
	//写难度值，难度值应当是推导出的，v2版本简化，把难度写成固定
	//0001//
	//targetStr := "0010000000000000000000000000000000000000000000000000000000000000"
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-bits)

	//bigIntTmp.SetString(targetStr, 16)
	pow.target = bigIntTmp
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	for {
		fmt.Printf("%x\n", hash)
		hash = sha256.Sum256(pow.prepareData(nonce))
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])
		if bigIntTmp.Cmp(pow.target) == -1 {
			fmt.Printf("wa kuang cheng gong ===== nonce: %d, hash: %x\n", nonce, hash)
			return hash[:], nonce
		}
		nonce++
	}
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulty),
		block.Data,
		uintToByte(nonce),
	}
	var data []byte
	data = bytes.Join(tmp, []byte{})
	return data
}
func (pow *ProofOfWork) IsValid() bool {
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var tmp big.Int
	tmp.SetBytes(hash[:])
	return tmp.Cmp(pow.target) == -1
}
