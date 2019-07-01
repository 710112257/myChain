package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}
//控制端
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		//发送者，接受者，发送金额
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

//打印链
func (cli *CLI) printChain() {

	bc := Newblockchain("")
	defer bc.db.Close()
	//调用区块链的迭代器，返回迭代器实例
	bci := bc.Iterator()
	//迭代输出区块链内容
	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
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

//添加块
func (cli *CLI) createBlockchain(address string) {
	bc := Newblockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}
//根据地址获取余额
func (cli *CLI) getBalance(address string) {
	//得到最新的区块链索引
	bc := Newblockchain(address)
	defer bc.db.Close()
	//初始化余额为0
	balance := 0
	//找到属于该地址的未花费集合
	UTXOs := bc.FindUTXO(address)
	//看一共还剩多少钱
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
//转账
func (cli *CLI) send(from, to string, amount int) {
	//得到最新的区块链索引
	bc := Newblockchain(from)
	defer bc.db.Close()
	//创建新的交易事务相当于DATA
	tx := NewUTXOTransaction(from, to, amount, bc)
	//添加入区块链中
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("Success!")
}



//命令行参数不正确报错
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS （查询余额）")
	fmt.Println("  createblockchain -address ADDRESS （创建区块链并发送Genesis区块奖励到地址）")
	fmt.Println("  printchain （打印区块链）")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT ( 转账)")
}
//验证命令行参数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}