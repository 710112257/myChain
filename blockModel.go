package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

//创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
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
func NewGenesisBlock()*Block{
	return NewBlock("这是创世区块",[]byte{})
}

//创建区块链实例
func Newblockchain() *Blockchain {
	var tip []byte
	//连接my.db
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//跟新数据
	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("blocksBucket"))
		//如果没有链则创建创世区块，否则找到索引
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte("blocksBucket"))
			if err != nil {
				fmt.Errorf(err.Error())
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}