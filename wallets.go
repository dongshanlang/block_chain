/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/5 10:28 AM
 * @Desc:
 */

package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//Wallets
type Wallets struct {
	WalletsMap map[string]*WalletKeyPair
}

//create Wallets
func NewWallets() *Wallets {
	ws := &Wallets{WalletsMap: make(map[string]*WalletKeyPair)}
	//1把所有的钱包从本地加载出来
	if !ws.LoadFromFile() {
		fmt.Println("load from file failed")
		return ws
	}

	//2返回实例

	return ws

}

const WalletName = "wallet.dat"

//这个wallets是对外的，WalletKeyPair是对内的
func (w *Wallets) CreateWallets() string {
	//调用NewWalletKeyPair
	wallet := NewWalletKeyPair()
	//将返回的WalletKeyPair添加到WalletMap中
	address := wallet.GetAddress()
	w.WalletsMap[address] = wallet
	//保存到本地
	res := w.SaveToFile()
	if !res {
		fmt.Println("save to file failed.")
		return ""
	}
	//返回新生成的key pair
	return address
}

//gob: type not registered for interface: elliptic.p256Curve
//gob interface类型的数据，需要告诉gob

func (w *Wallets) SaveToFile() bool {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256()) //注册interface类型数据
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(w)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(WalletName, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
	return true
}
func (w *Wallets) LoadFromFile() bool {
	if !IsFileExist(WalletName) {
		fmt.Println("wallet file not exist")
		return false
	}
	content, err := ioutil.ReadFile(WalletName)
	if err != nil {
		panic(err)
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(w)
	if err != nil {
		panic(err)
	}
	return true
}

func (w *Wallets) ListAddress() []string {
	var addresses []string
	for address, _ := range w.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
