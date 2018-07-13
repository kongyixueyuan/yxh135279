package BLC


// 创建创世区块
func (cli *CLI) yxh_createGenesisBlockchain(address string,nodeID string)  {

	blockchain := Yxh_CreateBlockchainWithGenesisBlock(address,nodeID)
	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.Yxh_ResetUTXOSet()
}

//blocks
//utxoTable