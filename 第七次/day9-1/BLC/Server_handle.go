package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
)

func handleVersion(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//Version
	//1. Version
	//2. BestHeight
	//3. 节点地址

	bestHeight := bc.Yxh_GetBestHeight()
	foreignerBestHeight := payload.Yxh_BestHeight

	if bestHeight > foreignerBestHeight {
		yxh_sendVersion(payload.Yxh_AddrFrom,bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		yxh_sendGetBlocks(payload.Yxh_AddrFrom)
	}


}

func handleAddr(request []byte,bc *Blockchain)  {




}

func handleGetblocks(request []byte,bc *Blockchain)  {


	var buff bytes.Buffer
	var payload GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()

	yxh_sendInv(payload.Yxh_AddrFrom, BLOCK_TYPE, blocks)


}

func handleGetData(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Yxh_Type == BLOCK_TYPE {

		block, err := bc.Yxh_GetBlock([]byte(payload.Yxh_Hash))
		if err != nil {
			return
		}

		yxh_sendBlock(payload.Yxh_AddrFrom, block)
	}

	if payload.Yxh_Type == "tx" {

	}


}

func handleBlock(request []byte,bc *Blockchain)  {
	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	block := payload.Yxh_Block

	bc.addBlock(block)

	//处理完一个后判断并继续下一个hash
	if len(yxh_transactionArray) > 0 {
		yxh_sendGetData(payload.Yxh_AddrFrom, BLOCK_TYPE, yxh_transactionArray[0])
		yxh_transactionArray = yxh_transactionArray[1:]
	} else {
		//数据库余额要重置
		utxoSet := &UTXOSet{bc}
		utxoSet.Yxh_ResetUTXOSet()
	}

}

func handleTx(request []byte,bc *Blockchain)  {

}


func handleInv(request []byte,bc *Blockchain)  {

	var buff bytes.Buffer
	var payload Inv

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	// Ivn 3000 block hashes [][]

	if payload.Yxh_Type == BLOCK_TYPE {

		yxh_transactionArray = payload.Yxh_Items

		blockHash := payload.Yxh_Items[0]
		yxh_sendGetData(payload.Yxh_AddrFrom, BLOCK_TYPE , blockHash)

		//判断长度，更新全局变量
		if len(payload.Yxh_Items) >= 1 {
			yxh_transactionArray = yxh_transactionArray[1:]
		}

	}

	if payload.Yxh_Type == TX_TYPE {

	}

}