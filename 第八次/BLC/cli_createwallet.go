package BLC

import "fmt"

func (cli *CLI) yxh_createWallet(nodeID string) {
	//wallet := Rwq_NewWallet()
	//address := wallet.Rwq_GetAddress()
	//fmt.Printf("钱包地址：%s\n",address)

	wallets, _ := Yxh_NewWallets(nodeID)
	address := wallets.Yxh_NewWallet().Yxh_GetAddress()
	wallets.Yxh_SaveToFile(nodeID)
	fmt.Printf("钱包地址：%s\n", address)

}
