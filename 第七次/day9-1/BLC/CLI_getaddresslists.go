package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI) yxh_addressLists(nodeID string)  {

	fmt.Println("打印所有的钱包地址:")

	wallets,_ := Yxh_NewWallets(nodeID)

	for address,_ := range wallets.Yxh_WalletsMap {

		fmt.Println(address)
	}
}