package BLC


// 创建创世区块
func (cli *CLI) yxh_createGenesisBlockchain(address string)  {

	blockchain := Yxh_CreateBlockchainWithGenesisBlock(address)
	//打开数据库后需要关闭链接
	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}
	//将交易保存到文件
	utxoSet.Yxh_ResetUTXOSet()
}
