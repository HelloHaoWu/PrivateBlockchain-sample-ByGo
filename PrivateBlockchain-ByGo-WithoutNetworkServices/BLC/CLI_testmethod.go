package BLC

import "fmt"

func (cli *CLI) TestMethod() {
	fmt.Println("正在将交易信息存储进本地数据库...")
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
	fmt.Println("存储成功.")
}
