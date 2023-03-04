package BLC

import "fmt"

// ↓用于查询指定地址余额
func (cli *CLI) getBalance(address string) {
	fmt.Println("地址: " + address)
	// ↓打开数据库, 获取最新状态的区块链
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	// ↓获取需要的余额
	//txOutputs := blockchain.UnUTXOs(address) // ← 这个功能添加测试的过程是尤为可贵的
	//for _, out := range txOutputs {
	//	fmt.Println(out)
	//}
	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.GetBalance(address)
	fmt.Printf("%s 一共有 %d 个Token.\n", address, amount)
}
