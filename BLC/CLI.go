package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
	"ketang/publicChain/BC/utils"
	"strconv"
)
type CLI struct {
	//BlockChain *BlockChain
}

func  (cli *CLI) Run()  {

	// 校验输入参数
	isValidArgs()
	//1.创建flagset命令对象
	createWalletCmd:=flag.NewFlagSet("createwallet",flag.ExitOnError)
	getAddresslistsCmd:=flag.NewFlagSet("getaddresslists",flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//2.设置命令后的参数对象
	flagcreateBlockChainData:=createBlockChainCmd.String("address",genesisCoinbaseData,"创世区块的数据")
	flagSendFromData := sendCmd.String("from", "", "发起转账者地址")
	flagSendToData := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountData := sendCmd.String("amount", "", "转账金额")

	flagGetBalanceData := getBalanceCmd.String("address", "", "要查询余额的账户")
	//3.解析
	switch os.Args[1] {
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getaddresslists":
		err := getAddresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		printUsage()
		os.Exit(1)

	}
	//4.根据终端输入的命令执行对应的功能
	if sendCmd.Parsed() {
		//fmt.Println("添加区块。。。",*flagAddBlockData)
		if *flagSendFromData == "" || *flagSendToData == "" || *flagSendAmountData == "" {
			fmt.Println("转账信息有误。。")
			printUsage()
			os.Exit(1)
		}
		//添加区块

		from := utils.JSONToArray(*flagSendFromData)     //[]string
		to := utils.JSONToArray(*flagSendToData)         //[]string
		amount := utils.JSONToArray(*flagSendAmountData) //[]string
		for i:=0;i<len(from);i++{
			if !IsValidAddress([]byte(from[i])) || !IsValidAddress([]byte(to[i])){
				fmt.Println("地址无效，无法转账。。")
				printUsage()
				os.Exit(1)
			}
			if a,_:=strconv.Atoi(amount[i]);a<=0 {
				fmt.Println("金额无效。")
				os.Exit(1)
			}
		}
		cli.Send(from, to, amount)
	}

	if printChainCmd.Parsed() {
		//fmt.Println("打印区块。。。")
		//cli.BlockChain.PrintChains()
		cli.PrintChains()
	}


	//添加创世区块的创建
	if createBlockChainCmd.Parsed() {
		if !IsValidAddress([]byte(*flagcreateBlockChainData)) {
			printUsage()
			os.Exit(1)
		}
		cli.CreateBlockChain(*flagcreateBlockChainData)
	}

	if getBalanceCmd.Parsed() {
		if  !IsValidAddress([]byte(*flagGetBalanceData)) {
			fmt.Println("查询地址有误。。")
			printUsage()
			os.Exit(1)
		}
		cli.GetBalance(*flagGetBalanceData)
	}
	if createWalletCmd.Parsed() {
		cli.CreateWallet()
	}
	if getAddresslistsCmd.Parsed(){
		cli.GetAddressLists()
	}

}

func printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\t createwallet -- 创建钱包")
	fmt.Println("\t getaddresslists -- 获取所有的钱包地址")
	fmt.Println("\t createblockchain -address DATA -- 创建创世区块")
	fmt.Println("\t send -from From -to To -amount Amount -- 转账交易")
	fmt.Println("\t printchain -- 打印区块")
	fmt.Println("\t getbalance -address Data -- 查询余额")
}

func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}










