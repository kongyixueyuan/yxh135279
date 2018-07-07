package BLC

import "fmt"

func (cli CLI) addresslists() {

	fmt.Println("打开印所有的钱包地址 ")

	wallets, _ := NewWallets()

	for address, _ := range wallets.Wallets {
		fmt.Println("address:",address)
	}


}