package BLC

import (
	"fmt"
	"os"
)

// ↓转账
func (cli *CLI) send(from []string, to []string, amount []string) {
	if dbExists() == false {
		fmt.Println("数据不存在.")
		os.Exit(1)
	}

	// ↓目的是获取blockchain上最新的block对象的索引, 但实际返回的是最新状态的blockchain
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close() // ←关闭数据库
	// 这里要补一个生成数字签名的方法
	// ↓ 在生成新的Block的过程中会验证数字签名是否合法
	blockchain.MineNewBlock(from, to, amount) // ← 创建t
}
