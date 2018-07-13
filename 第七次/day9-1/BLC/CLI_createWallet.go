package BLC

import "fmt"

func (cli *CLI) yxh_createWallet(nodeID string)  {

	wallets,_ := Yxh_NewWallets(nodeID)

	wallets.Yxh_CreateNewWallet(nodeID)

	fmt.Println(len(wallets.Yxh_WalletsMap))
}
