package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const genesisCoinbaseData = "This is zhe fist block!"
//创建区块
func NewBlock(Transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), Transactions, prevBlockHash, []byte{}, 0}
	//创建pow实例，数据存入实例中，并执行挖矿
	pow := NewProofOfWork(block)
	//此时得到达标的nonce和hash
	nonce, hash := pow.Run()
	//构建成功的区块并返回
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
//创建一个创世区块
func NewGenesisBlock(coinbase *Transaction)*Block{
	return NewBlock([]*Transaction{coinbase},[]byte{})
}

//创建区块链实例
func Newblockchain(address string) *Blockchain {
	var tip []byte
	//连接my.db
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//跟新数据
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		//如果没有链则创建创世区块，否则找到索引
		if b == nil {
			//生成了一笔交易
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			//生成创始区块
			genesis := NewGenesisBlock(cbtx)
			//连接到blocksBucket数据库
			b, err := tx.CreateBucket([]byte("blocksBucket"))
			if err != nil {
				fmt.Errorf(err.Error())
			}
			//哈希值为KEY，区块为value储存
			err = b.Put(genesis.Hash, genesis.Serialize())
			//并跟新最新KEY的索引
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	//最新的区块链索引
	bc := Blockchain{tip, db}

	return &bc
}