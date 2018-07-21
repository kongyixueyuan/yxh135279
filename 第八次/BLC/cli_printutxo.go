package BLC

func (cli *CLI) yxh_printutxo(nodeID string) {
	bc := Yxh_NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.Yxh_db.Close()
	UTXOSet.String()
}
