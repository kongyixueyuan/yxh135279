package BLC

import (
	"log"
	"fmt"
)

func (cli *CLI) yxh_getBalance(address string,nodeID string) {
	if !Yxh_ValidateAddress(address) {
		log.Panic("错误：地址无效")
	}

	bc := Yxh_NewBlockchain(nodeID)
	defer bc.Yxh_db.Close()
	UTXOSet := UTXOSet{bc}

	balance := UTXOSet.Yxh_GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}

func (cli *CLI) yxh_getBalanceAll(nodeID string) {
	wallets,err := Yxh_NewWallets(nodeID)
	if err!=nil{
		log.Panic(err)
	}
	balances := wallets.Yxh_GetBalanceAll(nodeID)
	for address,balance := range balances{
		fmt.Printf("地址:%s的余额为：%d\n", address, balance)
	}
}