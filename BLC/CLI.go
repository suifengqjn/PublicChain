package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
)
type CLI struct {
	//BlockChain *BlockChain
}

func  (cli *CLI) Run()  {

	// 校验输入参数
	isValidArgs()
	//1.创建flagset命令对象
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)

	//2.设置命令后的参数对象
	flagAddBlockData:=addBlockCmd.String("data","helloworld","区块的数据")
	//3.解析
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
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)

	}
	//4.根据终端输入的命令执行对应的功能
	if addBlockCmd.Parsed() {
		//fmt.Println("添加区块。。。",*flagAddBlockData)
		if *flagAddBlockData == ""{
			printUsage()
			os.Exit(1)
		}
		//添加区块
		cli.AddBlockToBlockChain(*flagAddBlockData)

	}

	if printChainCmd.Parsed() {
		//fmt.Println("打印区块。。。")
		//cli.BlockChain.PrintChains()
		cli.PrintChains()
	}

	//添加创世区块的创建
	if createBlockChainCmd.Parsed(){

		cli.CreateBlockChain(genesisCoinbaseData)
	}
}

func printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 添加区块")
	fmt.Println("\tprintchain -- 打印区块")
}

func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}


func (cli *CLI) PrintChains(){
	//cli.BlockChain.PrintChains()
	bc:=GetBlockChainObject()
	if bc == nil{
		fmt.Println("没有BlockChain，无法打印任何数据。。")
		os.Exit(1)
	}
	bc.PrintChains()
}


func (cli *CLI) AddBlockToBlockChain(data string){
	//cli.BlockChain.AddBlockToBlockChain(data)
	bc := GetBlockChainObject()
	if bc == nil{
		fmt.Println("没有BlockChain，无法添加新的区块。。")
		os.Exit(1)
	}
	bc.AddBlockToBlockChain([]byte(data))
}

func(cli *CLI) CreateBlockChain(data string){
	//fmt.Println("创世区块。。。")
	CreateBlockChainWithGenesisBlock([]byte(data))

}