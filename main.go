package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	targetBits=20
)
func main(){
	bc:=Newblockchain()
	for i:=0;i<=10;i++{
		str:=strconv.Itoa(i)
		startMine(str,bc)
	}
	fmt.Println(bc)

}
func startMine(i string,bc *Blockchain){
	startTime:=time.Now().UnixNano()
	var nonce=0
	for  {
		str:="这是第"+i+"个区块"
		hash:=Sha256(Sha256([]byte(str+strconv.Itoa(nonce))))
		hasgStr:=Getsha256string(hash)
		if strings.HasPrefix(hasgStr,strings.Repeat("0",targetBits)){
			bc.AddBlock(str)

			endTime:=float32((time.Now().UnixNano()-startTime)/1e9)
			fmt.Println("所用时间",endTime)
			fmt.Println()
			break
		}
		nonce=nonce+1
	}
}