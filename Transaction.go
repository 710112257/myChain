package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易，含有输入与输出
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}
//输入
type TXInput struct {
	Txid      []byte	//交易号
	Vout      int	//
	ScriptSig string	//签名
}
//输出
type TXOutput struct {
	Value        int  //值
	ScriptPubKey string //公钥脚本
}

const (
	subsidy  = 20
	)
//to 为地址 ，data 为内容，生成
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	//构建输入，
	txin := TXInput{[]byte{}, -1, data}
	//构建输出
	txout := TXOutput{subsidy, to}
	//构建交易
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
//解锁输出
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	//脚本值是否等于输入值
	return in.ScriptSig == unlockingData
}
//可用的解锁脚本
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
//发送交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//找到该用户的余额和能够花费的输出集合
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// 遍历能够花费的输出集合
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err!=nil{
			log.Fatal(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, from} //创建输入的实例
			inputs = append(inputs, input) //得到拥有所剩余额的实例输入
		}
	}

	// 创建新的输出实例，并且指定只有接收者能够使用
	outputs = append(outputs, TXOutput{amount, to})
	//如果余额大于转账金额
	if acc > amount {
		//创建新的输出实例，用户金额为减去转账金额，只有发送者能够使用
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}
	//创建打包好新的交易事务
	tx := Transaction{nil, inputs, outputs}
	//设置该事务的哈希ID
	tx.SetID()

	return &tx
}
//判断是否币基
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}
//以单个交易设置单个交易的哈希ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}