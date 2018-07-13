package BLC


func (cli *CLI) yxh_resetUTXOSet(nodeID string)  {

	blockchain := Yxh_BlockchainObject(nodeID)

	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.Yxh_ResetUTXOSet()

}
