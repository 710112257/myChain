package main

import (
	"crypto/sha256"
	"github.com/imroc/biu"
	"strings"
	"time"
)

//创建区块
func NewBlock(data string,prevBlockHash []byte) *Block {
	block:=&Block{time.Now().Unix(),[]byte(data),prevBlockHash,[]byte{}}
	block.SetHash()
	return block
}

//创建一个创世区块
func NewGenesisBlock()*Block{
	return NewBlock("这是创世区块",[]byte{})
}
//创建一个区块链对象
func Newblockchain() *Blockchain  {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
//哈希化
func Sha256(content []byte)[]byte{
	hash:=sha256.New()
	hash.Write(content)
	bs:=hash.Sum(nil)
	return bs
}
//把哈希进行二进制字符串查看
func Getsha256string(bs []byte)string{
	result:=biu.BytesToBinaryString(bs)
	result=result[1:len(result)-1]
	result=strings.Replace(result," ","",-1)
	return result
}