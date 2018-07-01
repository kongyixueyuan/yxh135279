package main

import (
	"fmt"
	"gostudy/blockchain/Block/BLC"
)

func main() {
	fmt.Println("Hello World")
	//block := BLC.CreateGenesisBlock("Genesis Block ...")
	//fmt.Println("block:",block)
	//fmt.Println("block.data:",string(block.Data))

	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	//fmt.Println("block.data:",string(blockchain.Blocks[0].Data))
	fmt.Println("blocks.len:",len(blockchain.Blocks))

	blockchain.AddBlockToBlockchain("Send 100RMB To zhangshang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,blockchain.Blocks[len(blockchain.Blocks)-1].Hash )
	blockchain.AddBlockToBlockchain("Send 200RMB To zhangshang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,blockchain.Blocks[len(blockchain.Blocks)-1].Hash )
	blockchain.AddBlockToBlockchain("Send 300RMB To zhangshang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,blockchain.Blocks[len(blockchain.Blocks)-1].Hash )
	blockchain.AddBlockToBlockchain("Send 400RMB To zhangshang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,blockchain.Blocks[len(blockchain.Blocks)-1].Hash )

	fmt.Println("blocks.len:",len(blockchain.Blocks))
	fmt.Println("block.data:",string(blockchain.Blocks[4].Data))


}