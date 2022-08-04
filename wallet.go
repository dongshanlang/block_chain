/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/4 6:33 PM
 * @Desc:
 */

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

func NewWallet() *WalletKeyPair {
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
