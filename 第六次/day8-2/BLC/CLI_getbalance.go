package BLC

import "fmt"

// 先用它去查询余额
func (cli *CLI) yxh_getBalance(address string)  {

	fmt.Println("地址：" + address)

	blockchain := Yxh_BlockchainObject()

	defer blockchain.Yxh_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.Yxh_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)

}
