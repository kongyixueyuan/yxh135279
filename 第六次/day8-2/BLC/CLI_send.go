package BLC

import (
	"fmt"
	"os"
)

// 转账
func (cli *CLI) yxh_send(from []string,to []string,amount []string)  {

	if Yxh_DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := Yxh_BlockchainObject()

	defer blockchain.Yxh_DB.Close()

	blockchain.Yxh_MineNewBlock(from,to,amount)

	utxoSet := &UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.Yxh_Update()

}

