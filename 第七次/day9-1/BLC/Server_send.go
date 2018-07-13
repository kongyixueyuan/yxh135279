package BLC

import (
	"fmt"
	"io"
	"bytes"
	"log"
	"net"
)


//COMMAND_VERSION
func yxh_sendVersion(toAddress string,bc *Blockchain)  {


	bestHeight := bc.Yxh_GetBestHeight()
	payload := yxh_gobEncode(Version{NODE_VERSION, bestHeight, yxh_nodeAddress})

	request := append(yxh_commandToBytes(COMMAND_VERSION), payload...)

	yxh_sendData(toAddress,request)


}



//COMMAND_GETBLOCKS
func yxh_sendGetBlocks(toAddress string)  {

	payload := yxh_gobEncode(GetBlocks{yxh_nodeAddress})

	request := append(yxh_commandToBytes(COMMAND_GETBLOCKS), payload...)

	yxh_sendData(toAddress,request)

}

// 主节点将自己的所有的区块hash发送给钱包节点
//COMMAND_BLOCK
//
func yxh_sendInv(toAddress string, kind string, hashes [][]byte) {

	payload := yxh_gobEncode(Inv{yxh_nodeAddress,kind,hashes})

	request := append(yxh_commandToBytes(COMMAND_INV), payload...)

	yxh_sendData(toAddress,request)

}



func yxh_sendGetData(toAddress string, kind string ,blockHash []byte) {

	payload := yxh_gobEncode(GetData{yxh_nodeAddress,kind,blockHash})

	request := append(yxh_commandToBytes(COMMAND_GETDATA), payload...)

	yxh_sendData(toAddress,request)
}


func yxh_sendData(to string,data []byte)  {

	fmt.Println("客户端向服务器发送数据......")
	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func yxh_sendBlock(toAddress string, block *Block) {

	payload := yxh_gobEncode(BlockData{toAddress,block})

	request := append(yxh_commandToBytes(COMMAND_BLOCK), payload...)

	yxh_sendData(toAddress,request)
}