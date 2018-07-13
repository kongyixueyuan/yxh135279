package BLC


func (cli *CLI) yxh_printchain(nodeID string)  {

	blockchain := Yxh_BlockchainObject(nodeID)

	defer blockchain.Yxh_DB.Close()

	blockchain.Yxh_Printchain()

}