package BLC

import (
	"fmt"
)

func (cli CLI) creteWallet() {

	wallets, _ := NewWallets()

	wallets.CreateNewWallet()

	fmt.Println(len(wallets.Wallets))
}