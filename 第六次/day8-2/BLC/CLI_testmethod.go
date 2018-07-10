package BLC

import "fmt"

// 重置方法
func (cli *CLI) Yxh_TestMethod()  {

	fmt.Println("TestMethod")

	blockchain := Yxh_BlockchainObject()

	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.Yxh_ResetUTXOSet()

	//fmt.Println(blockchain.FindUTXOMap())
}
