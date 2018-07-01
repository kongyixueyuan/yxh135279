package BLC

type Blockchain struct {
	Blocks []*Block //存储有序的区块
}

//1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {
	//创建创世区块
	genesisBlock := CreateGenesisBlock("Genesis Block ...")
	//返回区块链对象
	return &Blockchain{Blocks:[]*Block{genesisBlock}}

}

//2.新增普通区块并添加到区块链
func (blc *Blockchain)AddBlockToBlockchain(data string, height int64, preHash []byte) {
	//创建区块
	newBlock := NewBlock(data,height,preHash)
	//往区块链中添加区块
	blc.Blocks = append(blc.Blocks, newBlock)
}