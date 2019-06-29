package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}
//实例化共识机制
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//我们的目标是让一个目标在内存中占用少于256位
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}
//数据整理
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
//挖矿开始
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		//把不同的nonce存入获得一个完整的数据，再进行哈希
		data := pow.prepareData(nonce)
		//进行哈希加密
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		//当哈希达到预定值时，才算是合法的区块
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			//变化值只有nonce
			nonce++
		}
	}
	fmt.Print("\n\n")
	//返回达标的目标值和哈希
	return nonce, hash[:]
}
//验证块
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	//整理好数据格式
	data := pow.prepareData(pow.block.Nonce)
	//用已经上链的块生成哈希
	hash := sha256.Sum256(data)
	//是否与目标值一致
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}