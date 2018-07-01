package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}


//打印提示信息
func printUsage()  {

	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -data -- 交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")

}

//验证参数的的合法性
func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//添加区块
func (cli *CLI) addBlock(data string)  {
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
	blockchain.AddBlockToBlockchain(data)
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

//创建创世区块
func (cli *CLI) createGenesisBlockchain(data string)  {

	CreateBlockchainWithGenesisBlock(data)
}

//CLI入口函数
func (cli *CLI) Run()  {

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","http://yxh135279.org","交易数据......")

	flagCreateBlockchainWithData := createBlockchainCmd.String("data","Genesis block data......","创世区块交易数据......")

	//根据输入参数确定调用哪个入口
	switch os.Args[1] {
		case "addblock":
			err := addBlockCmd.Parse(os.Args[2:])
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
		default:
			printUsage()
			os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		//fmt.Println(*flagAddBlockData)
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithData)
	}

}