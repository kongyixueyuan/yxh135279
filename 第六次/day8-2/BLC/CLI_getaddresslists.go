package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI) yxh_addressLists()  {

	fmt.Println("打印所有的钱包地址:")

	wallets,_ := Yxh_NewWallets()

	for address,_ := range wallets.Yxh_WalletsMap {

		fmt.Println(address)
	}
}