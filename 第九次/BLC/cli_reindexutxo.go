package BLC

import "fmt"

func (cli *CLI) yxh_reindexUTXO(nodeID string)  {
	bc := Yxh_NewBlockchain(nodeID);
	defer bc.Yxh_db.Close()
	utxoset := UTXOSet{bc}
	utxoset.Yxh_Reset()
	fmt.Println("重建成功")
}
