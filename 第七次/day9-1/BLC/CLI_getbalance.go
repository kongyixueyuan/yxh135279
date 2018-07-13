package BLC

import "fmt"

// 先用它去查询余额
func (cli *CLI) yxh_getBalance(address string,nodeID string)  {

	fmt.Println("地址：" + address)

	// 获取某一个节点的blockchain对象
	blockchain := Yxh_BlockchainObject(nodeID)
	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.Yxh_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)

}
