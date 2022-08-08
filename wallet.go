/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/4 6:33 PM
 * @Desc:
 */

package main

import (
	"block_chain/base58"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

/*
1.创建结构体
type WalletKeyPair struct{}密钥对，保存公钥和私钥
2。给这个结构提供方法：
	1。GetAddress：私钥-》公钥-》地址
3。地址生成规则
*/

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte //可以将公钥的X、Y拼接后传输，在对端再进行切割还原，好处是可以方便后面的编码
}

func NewWalletKeyPair() *WalletKeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	publicKeyRaw := privateKey.PublicKey
	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
func (w *WalletKeyPair) GetAddress() string {
	publicHash := HashPublicKey(w.PublicKey)

	version := 0x00
	//21个字节
	payload := append([]byte{byte(version)}, publicHash...)

	checksum := CheckSum(payload) //4字节校验码

	payload = append(payload, checksum...) //一共25个字节

	address := base58.Encode(payload)
	return address
}
func HashPublicKey(publicKey []byte) []byte {
	hash := sha256.Sum256(publicKey)
	//创建hash160对象
	//向hash160中write数据
	//做hash运算
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])
	if err != nil {
		panic(err)
	}
	publicHash := rip160Hasher.Sum(nil)
	return publicHash
}
func CheckSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	checksum := second[0:4] //4字节校验码
	return checksum
}
func IsValidAddress(address string) bool {
	//1.将输入的地址进行解码得到25字节
	//2。去除前21字节，运行check sum函数，得到checksum1
	//3。取出后4字节，得到checksum2
	//4。比较两个checksum，如果地址相同有效，否则无效
	decodeInfo := base58.Decode(address)
	if len(decodeInfo) != 25 { //25是比特币的地址长度，固定值
		return false
	}
	payload := decodeInfo[:len(decodeInfo)-4]
	checkSum1 := CheckSum(payload)
	checkSum2 := decodeInfo[len(decodeInfo)-4:]
	if bytes.Equal(checkSum2, checkSum1) {
		return true
	}
	return false
}
