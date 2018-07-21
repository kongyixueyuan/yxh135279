package BLC

import "log"

func (cli *CLI) yxh_createblockchain(address string,nodeID string)  {
	//验证地址是否有效
	if !Yxh_ValidateAddress(address){
		log.Panic("地址无效")
	}
	bc := Yxh_CreateBlockchain(address,nodeID)
	defer bc.Yxh_db.Close()

	// 生成UTXOSet数据库
	UTXOSet := UTXOSet{bc}
	UTXOSet.Yxh_Reindex()
}
