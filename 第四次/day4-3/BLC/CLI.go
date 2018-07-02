package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}


//打印提示信息
//main.exe send -from "[\"yxh\"]" -to "[\"yxh1\"]" -amount "[\"10\"]"
//main.exe send -from "[\"yxh\",\"yxh1\"]" -to "[\"yxh1\","\yxh2\"]" -amount "[\"10\","\2\"]"
func printUsage()  {

	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from ADDRESS -to ADDRESS -amount VALUE -- 转帐.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address ADDRESS -- 查询余额.")

}

//验证参数的的合法性
func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//添加区块
func (cli *CLI) send(from []string, to []string, amount []string)  {
	//判断数据库是否存在
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	//通过从数据库中查找最新的区块hash
	blockchain := BlockchainObject()
	//最后新关闭数据库链接
	defer blockchain.DB.Close()
	//添加新区块
	//组织Transaction


	blockchain.MineNewBlock(from, to, amount)
}

//打印所有区块信息
func (cli *CLI) printchain()  {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	//通过从数据库中查找最新的区块hash
	blockchain := BlockchainObject()
	//最后新关闭数据库链接
	defer blockchain.DB.Close()
	//打印区块信息
	blockchain.Printchain()

}

func getbalance(address string) {

	//通过从数据库中查找最新的区块hash
	blockchain := BlockchainObject()

	//最后新关闭数据库链接
	defer blockchain.DB.Close()
	//打印区块信息
	amount := blockchain.GetBalance(address)

	fmt.Println("%s一共有%dToken", address, amount)

}

//创建创世区块
func (cli *CLI) createGenesisBlockchain(data string)  {

	CreateBlockchainWithGenesisBlock(data)
}

// 先用它去查询余额
func (cli *CLI) getBalance(address string)  {

	fmt.Println("地址：" + address)

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)


}

//CLI入口函数
func (cli *CLI) Run()  {

	isValidArgs()

	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")


	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")

	flagbalance := getbalanceCmd.String("address","","查询金额......")


	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
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

	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}



		//fmt.Println(*flagAddBlockData)
		//cli.addBlock([]*Transaction{})

		//fmt.Println(*flagFrom)
		//fmt.Println(*flagTo)
		//fmt.Println(*flagAmount)

		fmt.Println(JSONToArray(*flagFrom))
		fmt.Println(JSONToArray(*flagTo))
		fmt.Println(JSONToArray(*flagAmount))

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

	if printChainCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {

		if *flagbalance == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*flagbalance)
	}

}