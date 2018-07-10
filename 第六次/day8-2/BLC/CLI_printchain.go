package BLC

import (
	"fmt"
	"os"
)

// 打印区块信息
func (cli *CLI) yxh_printchain()  {

	if Yxh_DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := Yxh_BlockchainObject()

	defer blockchain.Yxh_DB.Close()

	blockchain.Yxh_Printchain()

}