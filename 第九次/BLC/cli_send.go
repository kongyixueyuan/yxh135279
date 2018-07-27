package BLC

func (cli *CLI) yxh_send(from []string, to []string, amount []string,nodeID string, mineNow bool) {
	bc := Yxh_NewBlockchain(nodeID)
	defer bc.Yxh_db.Close()
	bc.MineNewBlock(from, to, amount,nodeID, mineNow)
}