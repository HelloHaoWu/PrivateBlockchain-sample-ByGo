package main

// ↓每次复制, go.mod和main.go都要改
import (
	"Part54-UTXOSet-GetBalance-P93/BLC"
)

func main() {
	//var blockchain *BLC.Blockchain
	//// 创建包含创世区块的区块链
	//if len(os.Args) == 1 {
	//	blockchain_create := BLC.CreateBlockchainWithGenesisBlock()
	//	blockchain = blockchain_create
	//}

	//blockchain := BLC.CreateBlockchainWithGenesisBlock(os.Args)

	//defer blockchain.DB.Close()
	// ↑ 不能关, 得随时更新随时打开

	//// 添加新区块
	//blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang")
	//blockchain.AddBlockToBlockchain("Send 200RMB To changjingkong")
	//blockchain.AddBlockToBlockchain("Send 300RMB To juncheng")
	//blockchain.AddBlockToBlockchain("Send 50RMB To haolin")
	//
	//// 遍历区块链
	//blockchain.Printchain()

	cli := &BLC.CLI{}
	// ↑ 这为什么转指针?? → 因为输入就是指针
	cli.Run()
}
