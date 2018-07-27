package BLC

import (
	"log"
	"fmt"
)

func (cli *CLI) yxh_listAddrsss(nodeID string)  {
	wallets,err := Yxh_NewWallets(nodeID)

	if err!=nil{
		log.Panic(err)
	}
	addresses := wallets.Yxh_GetAddresses()

	for _,address := range addresses{
		fmt.Printf("%s\n",address)
	}
}
