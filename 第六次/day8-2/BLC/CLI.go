package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}

func yxh_printUsage()  {

	fmt.Println("Usage:")

	fmt.Println("\taddresslists -- 输出所有钱包地址.")
	fmt.Println("\tcreatewallet -- 创建钱包.")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")
	fmt.Println("\ttest -- 测试.")

}

//校验输入是否合法
func yxh_isValidArgs()  {
	if len(os.Args) < 2 {
		yxh_printUsage()
		os.Exit(1)
	}
}

//开始运行
func (cli *CLI) Yxh_Run()  {

	yxh_isValidArgs()

	testCmd := flag.NewFlagSet("test",flag.ExitOnError)
	addresslistsCmd := flag.NewFlagSet("addresslists",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address","","要查询某一个账号的余额.......")

	switch os.Args[1] {
		case "send":
			err := sendBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "test":
			err := testCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "addresslists":
			err := addresslistsCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "getbalance":
			err := getbalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createwallet":
			err := createWalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			yxh_printUsage()
			os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			yxh_printUsage()
			os.Exit(1)
		}

		from := Yxh_JSONToArray(*flagFrom)
		to := Yxh_JSONToArray(*flagTo)

		for index,fromAdress := range from {
			if Yxh_IsValidForAdress([]byte(fromAdress)) == false || Yxh_IsValidForAdress([]byte(to[index])) == false {
				fmt.Printf("地址无效......")
				yxh_printUsage()
				os.Exit(1)
			}
		}

		amount := Yxh_JSONToArray(*flagAmount)
		cli.yxh_send(from,to,amount)
	}

	if printChainCmd.Parsed() {
		//fmt.Println("输出所有区块的数据........")
		cli.yxh_printchain()
	}

	if testCmd.Parsed() {
		fmt.Println("测试....")
		cli.Yxh_TestMethod()
	}

	if addresslistsCmd.Parsed() {
		//fmt.Println("输出所有区块的数据........")
		cli.yxh_addressLists()
	}


	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.yxh_createWallet()
	}

	if createBlockchainCmd.Parsed() {

		if Yxh_IsValidForAdress([]byte(*flagCreateBlockchainWithAddress)) == false {
			fmt.Println("地址无效....")
			yxh_printUsage()
			os.Exit(1)
		}

		cli.yxh_createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {

		if Yxh_IsValidForAdress([]byte(*getbalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			yxh_printUsage()
			os.Exit(1)
		}

		cli.yxh_getBalance(*getbalanceWithAdress)
	}

}