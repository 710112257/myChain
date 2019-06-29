package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}
//控制端
func (cli *CLI) Run() {
	//验证命令参数
	cli.validateArgs()

	//我们使用标准的标志包来解析命令行参数
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			fmt.Errorf(err.Error())
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err!=nil{
			fmt.Errorf(err.Error())
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	//如果指令为addBlock执行这个
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}
	//如果指令为printChain执行这个
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
//添加块
func (cli *CLI) addBlock(data string) {
	//添加新块
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}
//打印链
func (cli *CLI) printChain() {
	//调用区块链的迭代器，返回迭代器实例
	bci := cli.bc.Iterator()

	//迭代输出区块链内容
	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		//创建pow实例，并检查块是否合法
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
//命令行参数不正确报错
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}
//验证命令行参数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}