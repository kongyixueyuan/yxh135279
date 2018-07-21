package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

const walletFile  = "wallet_%s.dat"

type Wallets struct {
	Yxh_Wallets map[string]*Wallet
}

// 生成新的钱包
// 从数据库中读取，如果不存在
func Yxh_NewWallets(nodeID string)(*Wallets,error)  {
	wallets := Wallets{}
	wallets.Yxh_Wallets = make(map[string]*Wallet)

	err := wallets.Yxh_LoadFromFile(nodeID)

	return &wallets,err
}
// 生成新的钱包地址列表
func (ws *Wallets) Yxh_NewWallet() *Wallet {
	wallet := Yxh_NewWallet()
	address := wallet.Yxh_GetAddress()
	ws.Yxh_Wallets[string(address)] = wallet
	return wallet
}
// 获取钱包地址
func (ws *Wallets) Yxh_GetAddresses()[]string  {
	var addresses []string
	for address := range ws.Yxh_Wallets{
		addresses = append(addresses,address)
	}
	return addresses
}

// 根据地址获取钱包的详细信息
func (ws Wallets) Yxh_GetWallet(address string) Wallet {
	return *ws.Yxh_Wallets[address]
}

// 从数据库中读取钱包列表
func (ws *Wallets)Yxh_LoadFromFile(nodeID string) error  {
	 walletFile := fmt.Sprintf(walletFile, nodeID)
	 if _,err := os.Stat(walletFile) ; os.IsNotExist(err){
	 	return err
	 }

	 fileContent ,err := ioutil.ReadFile(walletFile)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 var wallets Wallets
	 gob.Register(elliptic.P256())
	 decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	 err = decoder.Decode(&wallets)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 ws.Yxh_Wallets = wallets.Yxh_Wallets

	 return nil
}

// 将钱包存到数据库中
func (ws *Wallets)Yxh_SaveToFile(nodeID string)  {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err !=nil{
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err !=nil{
		log.Panic(err)
	}
}
// 打印所有钱包的余额
func (ws *Wallets) Yxh_GetBalanceAll(nodeID string) map[string]int {
	addresses := ws.Yxh_GetAddresses()
	bc := Yxh_NewBlockchain(nodeID)
	defer bc.Yxh_db.Close()
	UTXOSet := UTXOSet{bc}

	result := make(map[string]int)
	for _,address := range addresses{
		if !Yxh_ValidateAddress(address) {
			result[address] = -1
		}
		balance := UTXOSet.Yxh_GetBalance(address)
		result[address] = balance
	}
	return result
}