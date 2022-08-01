/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 5:45 PM
 * @Desc:
 */

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Version       uint64
	PrevBlockHash []byte
	MerKleRoot    []byte //
	TimeStamp     uint64 //秒
	Difficulty    uint64 //难度
	Nonce         uint64 //随机数，挖矿即要计算的随机数
	Hash          []byte // current block hash
	Data          []byte
}

//create block,
func NewBlock(data string, preBlockHash []byte) *Block {
	block := Block{
		Version:       1.0,
		PrevBlockHash: preBlockHash,
		MerKleRoot:    nil,
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulty:    10,
		//Nonce:         10,
		Hash: nil,
		Data: []byte(data),
	}
	if string(preBlockHash) == string(FirstBlock) {
		block.PrevBlockHash = nil
	}
	block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

func (block *Block) SetHash() {
	var data []byte
	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulty),
		block.Data,
		uintToByte(block.Nonce),
	}
	data = bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
func (block *Block) Deserialize(data []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		panic(err)
	}
	return &b
}

//func (block *Block) Serialize() []byte {
//	var buffer bytes.Buffer
//	encoder := gob.NewEncoder(&buffer)
//	err := encoder.Encode(block)
//	if err != nil {
//		panic(err)
//	}
//	return buffer.Bytes()
//}
func Deserialize(data []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		panic(err)
	}
	return &b
}