package main

const (
	targetBits=20
)
func main(){
	bc := Newblockchain()
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
