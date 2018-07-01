package main

import (
	"fmt"
	"gostudy/blockchain/day2-demo1/BLC"
)

func main() {
	fmt.Println("Hello World")

	genesisBlock := BLC.CreateGenesisBlock("Genesis Block ...")
	fmt.Println(genesisBlock.Nonce)
	fmt.Println(genesisBlock)
	genblockBytes := genesisBlock.Serialize()
	fmt.Println(genblockBytes)
	block := BLC.DeserializeBlock(genblockBytes)
	fmt.Println(block)
	fmt.Println(block.Nonce)
}
