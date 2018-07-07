package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"io/ioutil"
	"crypto/elliptic"
	"os"
)
//钱包文件名
const walletFile = "wallets.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}


// 创建钱包集合
func NewWallets() (*Wallets, error) {

	var wallets Wallets

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets = Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return &wallets, err
	}

	//从文件中读取钱包所有的地址
	fileContent, err := ioutil.ReadFile(walletFile)

	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallets)

	if err != nil {
		log.Panic(err)
	}

	return &wallets, err
}

// 创建一个新钱包
func (w *Wallets) CreateNewWallet()  {

	wallet := NewWallet()
	fmt.Printf("Address：%s\n",wallet.GetAddress())
	w.Wallets[string(wallet.GetAddress())] = wallet
	w.saveWallets()
}

//保存钱包
func (w *Wallets) saveWallets() {
	var content bytes.Buffer

	//注册的目的，是为了可以将任何类型序列化
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(w)

	if err != nil {

		log.Panic(err)

	}

	//写文件 ，没有则新建，有则重写
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)

	if err != nil {

		log.Panic(err)

	}

}


