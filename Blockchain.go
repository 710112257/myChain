package main

import "fmt"

type Blockchain struct {
	Block []*Block
}

//增加新块
func (bc *Blockchain)AddBlock(data string)  {
	preblockhash:=bc.Block[len(bc.Block)-1].Hash
	block:=NewBlock(data,preblockhash)
	fmt.Printf("时间戳：%x\n",block.Timestamp)
	fmt.Printf("数据：%s\n",block.Data)
	fmt.Printf("当前哈希：%x\n",block.Hash)
	fmt.Printf("上一个哈希：%x\n",block.PrevBlockHash)
	bc.Block=append(bc.Block,block)
}