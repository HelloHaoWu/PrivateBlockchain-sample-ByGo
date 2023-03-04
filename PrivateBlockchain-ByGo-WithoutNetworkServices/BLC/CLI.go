package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 这个文件是命令行工具部分(或者说实现命令行工具的相关代码)
type CLI struct{}

func printUsage() {
	fmt.Println("Usage: (--及其后面为提示, 使用时不需要输入--再输入所需参数, 直接在--位置输入参数就好)")
	fmt.Println("\tCreateWallets -- 创建钱包")
	fmt.Println("\tGetAddresslists -- 输出所有钱包地址")
	fmt.Println("\tCreateBlockchain -address -- 交易数据")
	fmt.Println("\tsend -from From -to To -amount AMOUNT -- 交易明细")
	fmt.Println("\tgetbalance -address -- 获取对应账户信息") // ← 加新功能改动点1
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tupdate -- 手动更新交易数据库(每次交易后会自动更新)")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//func (cli *CLI) addblock(txs []*Transaction) {
//	if dbExists() == false {
//		fmt.Println("数据不存在.")
//		os.Exit(1)
//	}
//	blockchain := GetBlockchainObject()
//	defer blockchain.DB.Close()
//	blockchain.AddBlockToBlockchain(txs)
//}

// ↓这里原始传的是一个已创建好的空交易
//func (cli *CLI) createGenesisBlockchain(txs []*Transaction) {
//	//fmt.Println(data)
//	CreateBlockchainWithGenesisBlock(txs)
//}

func (cli *CLI) Run() {
	isValidArgs()

	testCmd := flag.NewFlagSet("test", flag.ExitOnError) // 测试用, 之后要删除
	createWalletCmd := flag.NewFlagSet("CreateWallets", flag.ExitOnError)
	addresslistsCmd := flag.NewFlagSet("GetAddresslists", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("CreateBlockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError) // ← 加新功能改动点2

	flagFrom := sendBlockCmd.String("from", "", "转账源地址")
	flagTo := sendBlockCmd.String("to", "", "转账目的地地址")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额")

	flagcreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址...")
	// ↑ ./main -data 调用上述代码
	getbalanceWithAddress := getbalanceCmd.String("address", "", "要查询某一账号的余额") // ← 加新功能改动点3

	switch os.Args[1] {
	// ↑ 直接./main不输入会报错
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "test": // 测试用, 之后删除
		err := testCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "CreateBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance": // ← 加新功能改动点4
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "CreateWallets":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "GetAddresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1) // 默认1, 1退出
	}
	if createBlockchainCmd.Parsed() {
		if ValidateAddress([]byte(*flagcreateBlockchainWithAddress)) == false {
			fmt.Println("EOF: 地址无效.")
			printUsage()
			os.Exit(1) // ← 程序退出
		}
		cli.createGenesisBlockchain(*flagcreateBlockchainWithAddress)
	}
	if sendBlockCmd.Parsed() {
		// ↓严格来说, 还需要判断输入的三个部分的数量相等
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagAddBlockData) //-data后面的内容
		//fmt.Printf("flag: %s\n", *flagAddBlockData)
		//cli.addblock([]*Transaction{})
		//fmt.Println(*flagFrom)
		//fmt.Println(*flagTo)
		//fmt.Println(*flagAmount)
		//
		//fmt.Println(JsonToArray(*flagFrom))
		//fmt.Println(JsonToArray(*flagTo))
		//fmt.Println(JsonToArray(*flagAmount))

		from := JsonToArray(*flagFrom)
		to := JsonToArray(*flagTo)

		// ↓判断地址是否有效
		for index, fromAddress := range from {
			if ValidateAddress([]byte(fromAddress)) == false || ValidateAddress([]byte(to[index])) == false {
				fmt.Printf("地址无效或地址不合法.")
				printUsage()
				os.Exit(1) // ← 退出
			}
		}

		amount := JsonToArray(*flagAmount)
		cli.send(from, to, amount)
		cli.TestMethod() // 更新数据库
	}
	if printChainCmd.Parsed() {
		cli.printchain()
	}
	if testCmd.Parsed() { // 测试用, 之后删除
		cli.TestMethod()
	}
	if createWalletCmd.Parsed() {
		// ↓ 创建钱包
		cli.createWallet()
	}
	if addresslistsCmd.Parsed() {
		cli.addressLists()
	}
	if getbalanceCmd.Parsed() {
		if ValidateAddress([]byte(*getbalanceWithAddress)) == false {
			fmt.Println("EOF: 地址无效.")
			printUsage()
			os.Exit(1) // ← 程序退出
		}
		cli.getBalance(*getbalanceWithAddress)
	}
}
