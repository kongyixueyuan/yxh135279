package BLC

import "fmt"

// 转账
func (cli *CLI) yxh_send(from []string,to []string,amount []string,nodeID string,mineNow bool)  {


	blockchain := Yxh_BlockchainObject(nodeID)
	defer blockchain.Yxh_DB.Close()

	if mineNow {
		blockchain.Yxh_MineNewBlock(from,to,amount,nodeID)

		utxoSet := &UTXOSet{blockchain}

		//转账成功以后，需要更新一下
		utxoSet.Yxh_Update()
	} else {
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}



}

