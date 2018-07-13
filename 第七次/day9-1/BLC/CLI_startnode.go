package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) yxh_startNode(nodeID string,minerAdd string)  {

	// 启动服务器

	if minerAdd == "" || Yxh_IsValidForAdress([]byte(minerAdd))  {
		//  启动服务器
		fmt.Printf("启动服务器:localhost:%s\n",nodeID)
		yxh_startServer(nodeID,minerAdd)

	} else {

		fmt.Println("指定的地址无效....")
		os.Exit(0)
	}

}