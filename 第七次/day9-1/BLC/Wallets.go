package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"log"
	"os"
)

const yxh_walletFile  = "Wallets_%s.dat"

type Wallets struct {
	Yxh_WalletsMap map[string]*Wallet
}



// 创建钱包集合
func Yxh_NewWallets(nodeID string) (*Wallets,error){

	walletFile := fmt.Sprintf(yxh_walletFile,nodeID)

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Yxh_WalletsMap = make(map[string]*Wallet)
		return wallets,err
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

	return &wallets,nil
}





// 创建一个新钱包
func (w *Wallets) Yxh_CreateNewWallet(nodeID string)  {

	wallet := Yxh_NewWallet()
	fmt.Printf("Address：%s\n",wallet.Yxh_GetAddress())
	w.Yxh_WalletsMap[string(wallet.Yxh_GetAddress())] = wallet
	w.Yxh_SaveWallets(nodeID)
}

// 将钱包信息写入到文件
func (w *Wallets) Yxh_SaveWallets(nodeID string)  {


	walletFile := fmt.Sprintf(yxh_walletFile,nodeID)

	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}


}

