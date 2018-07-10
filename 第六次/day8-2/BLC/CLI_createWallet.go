package BLC

import "fmt"

//创建钱包
func (cli *CLI) yxh_createWallet()  {

	//钱包集
	wallets,_ := Yxh_NewWallets()
	//创建钱包
	wallets.Yxh_CreateNewWallet()

	fmt.Println(len(wallets.Yxh_WalletsMap))
}
