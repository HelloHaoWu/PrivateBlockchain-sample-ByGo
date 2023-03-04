package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "Wallets.dat"

type Wallets struct { // ← 无序, 以字典形式进行存储
	Wallets map[string]*Wallet // ← 地址映射Wallet
}

// ↓ 创建自己的Wallets
func NewWallets() (*Wallets, error) {
	// ↓ 判断是否存在, 不存在则创建
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets, err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets, nil
}

// ↓ 创建一个新钱包(调用方法就会)
func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Printf("Address of New Wallet: %s\n", wallet.GetAddress())
	w.Wallets[string(wallet.GetAddress())] = wallet
}

// ↓ 将钱包信息写入文件
func (w *Wallets) SaveWallets() {
	var content bytes.Buffer

	// 注册的目的是 → 可以序列化任何类型(接口之类的)
	gob.Register(elliptic.P256()) // ← 这条代码在go 1.19会出现"panic: gob: type elliptic.p256Curve has no exported fields"的bug
	// ↑ 解决方法是版本降级为1.18

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	// ↓ 将序列化以后的数据写入文件, 原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
