package BLC

// ↓创建创世区块
func (cli *CLI) createGenesisBlockchain(address string) {
	//fmt.Println(data)
	blockchain := CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
}
